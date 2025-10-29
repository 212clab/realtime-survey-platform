package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"user-service/internal/models"
	"user-service/internal/storage"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_super_secret_key")

// API는 HTTP 핸들러를 위한 구조체입니다.
type API struct {
	Store *storage.Store
}

// NewAPI는 새로운 API 인스턴스를 생성합니다.
func NewAPI(store *storage.Store) *API {
	return &API{Store: store}
}

// SignupHandler는 회원가입 요청을 처리합니다.
func (a *API) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := a.Store.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// LoginHandler는 로그인 요청을 처리하고 JWT를 발급합니다.
func (a *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if creds.Username == nil || creds.Password == nil {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := a.Store.GetUserByUsername(*creds.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(*creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// GoogleLoginHandler는 Google OAuth 콜백을 처리합니다.
func (a *API) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := "http://localhost:3000/api/auth/callback/google"

	tokenURL := "https://oauth2.googleapis.com/token"
	reqBody, _ := json.Marshal(map[string]string{
		"code":          code,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	})

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get google access token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tokenData)
	accessToken := tokenData["access_token"].(string)

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get google user info", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()
	
	var googleUser models.GoogleUserResponse
	json.NewDecoder(userResp.Body).Decode(&googleUser)

	userID, err := a.Store.FindOrCreateUserByGoogle(&googleUser)
	if err != nil {
		http.Error(w, "Failed to process google user data", http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// HealthCheckHandler는 서비스 상태를 확인합니다.
func (a *API) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}