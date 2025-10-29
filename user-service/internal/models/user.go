package models

// User는 사용자 데이터 구조체입니다.
type User struct {
	Username *string `json:"username"`
    Password *string `json:"password"`
    ID *int64  `json:"id,omitempty"` // 👈 필드 이름을 명확하게!
}


type GoogleUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}