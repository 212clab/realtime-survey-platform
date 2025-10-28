package models

// User는 사용자 데이터 구조체입니다.
type User struct {
	Username *string `json:"username"`
    Password *string `json:"password"`
    GithubID *int64  `json:"github_id,omitempty"` // 👈 필드 이름을 명확하게!
}

// GitHubUserResponse는 GitHub API로부터 받아오는 사용자 정보 구조체입니다.
type GitHubUserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email *string `json:"email"` // 이메일은 비공개일 수 있으므로 포인터로 처리
}

type GoogleUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}