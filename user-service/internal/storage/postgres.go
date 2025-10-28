package storage

import (
	"database/sql"
	"errors"                       // errors 패키지 import
	"user-service/internal/models" // 우리가 만든 models 패키지를 import

	"golang.org/x/crypto/bcrypt"
)

// Store는 데이터베이스 작업을 위한 구조체입니다.
type Store struct {
	DB *sql.DB
}

// NewStore는 새로운 Store 인스턴스를 생성합니다.
func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

// CreateUser는 새로운 사용자를 데이터베이스에 생성합니다.
func (s *Store) CreateUser(user *models.User) error {
	// 1. 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. DB에 사용자 정보 삽입
	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = s.DB.Exec(sqlStatement, user.Username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
// UserForLogin은 로그인 검증을 위해 DB에서 가져올 데이터 구조입니다.
type UserForLogin struct {
	ID             int
	Username       string
	HashedPassword string
}

// GetUserByUsername은 username으로 사용자를 찾아 비밀번호 해시와 함께 반환합니다.
func (s *Store) GetUserByUsername(username string) (*UserForLogin, error) {
	var user UserForLogin
	sqlStatement := `SELECT id, username, password FROM users WHERE username=$1`
	err := s.DB.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// 사용자가 없는 것은 에러가 아니라, 그냥 결과가 없는 것
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}