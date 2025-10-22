package models

// User는 사용자 데이터 구조체입니다.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}