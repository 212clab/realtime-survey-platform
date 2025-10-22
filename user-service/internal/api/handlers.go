package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/storage" // storage 패키지도 import
)

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