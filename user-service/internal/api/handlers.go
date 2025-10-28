package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time" // time 패키지 import
	"user-service/internal/models"
	"user-service/internal/storage"

	jwt "github.com/golang-jwt/jwt/v5" // jwt 라이브러리 import
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

	// DB 작업은 storage 계층에 위임
	if err := a.Store.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// HealthCheckHandler는 서비스 상태를 확인합니다.
func (a *API) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// 이 부분은 storage를 통해 DB 상태를 확인하도록 개선할 수도 있습니다.
	fmt.Fprintf(w, "OK") // 이 코드 때문에 fmt 패키지가 필요합니다.
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
	// 1. DB에서 사용자 정보 가져오기
	user, err := a.Store.GetUserByUsername(*creds.Username) // 👈 creds.Username -> *creds.Username
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 2. 입력된 비밀번호와 DB의 해시값 비교
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(*creds.Password)); err != nil { // 👈 creds.Password -> *creds.Password
		// 비밀번호 불일치
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. JWT 생성
	expirationTime := time.Now().Add(24 * time.Hour) // 토큰 유효기간: 24시간
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID), // 토큰 주체: 사용자 ID
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 4. 토큰을 응답으로 전송
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

type GitHubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type GitHubUserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}


// GithubLoginHandler는 GitHub OAuth 콜백을 처리합니다.
func (a *API) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 프론트엔드로부터 임시 허가증(code)을 받습니다.
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	// 2. 허가증을 Access Token으로 교환합니다.
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	})
	
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResp GitHubAccessTokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	// 3. Access Token으로 사용자 정보를 가져옵니다.
	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	
	userResp, err := client.Do(userReq)
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var githubUser GitHubUserResponse
	json.NewDecoder(userResp.Body).Decode(&githubUser)

	// 4. 받아온 사용자 정보로 우리 DB에 사용자를 생성하거나 찾습니다.
	//    (이 부분은 storage 계층에 새로운 함수를 만들어 처리해야 합니다.)
	//    예: user, err := a.Store.FindOrCreateUserByGithub(githubUser)
	//    ...

	// 5. 우리 서비스의 JWT를 발급합니다.
	//    (기존 LoginHandler의 JWT 생성 로직을 재사용)
	//    ...

	// 6. 최종적으로 JWT를 프론트엔드에 전달합니다.
	w.Write([]byte("GitHub Login Success! (JWT 발급 로직 추가 필요)"))
}