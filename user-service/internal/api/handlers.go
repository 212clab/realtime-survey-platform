package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	// 1. DB에서 사용자 정보 가져오기
	user, err := a.Store.GetUserByUsername(creds.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 2. 입력된 비밀번호와 DB의 해시값 비교
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(creds.Password)); err != nil {
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