package storage

import (
	"database/sql"
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