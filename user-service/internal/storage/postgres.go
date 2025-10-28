package storage

import (
	"database/sql" // sql 패키지를 사용합니다.
	"errors"
	"user-service/internal/models"

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
	if user.Username == nil || user.Password == nil {
		return errors.New("username and password are required for standard signup")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = s.DB.Exec(sqlStatement, *user.Username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

// UserForLogin은 로그인 검증을 위해 DB에서 가져올 데이터 구조입니다.
type UserForLogin struct {
	ID             int
	Username       string
	HashedPassword string // 이 구조체는 그대로 string을 유지합니다.
}

// GetUserByUsername은 username으로 사용자를 찾아 비밀번호 해시와 함께 반환합니다.
func (s *Store) GetUserByUsername(username string) (*UserForLogin, error) {
	var user UserForLogin
	// 👇 password를 받을 변수를 sql.NullString으로 선언합니다.
	var hashedPassword sql.NullString

	sqlStatement := `SELECT id, username, password FROM users WHERE username=$1`
	// 👇 Scan할 때도 hashedPassword 변수를 사용합니다.
	err := s.DB.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 👇 password가 NULL인지(유효하지 않은지) 확인합니다.
	if !hashedPassword.Valid {
		// 사용자는 있지만 비밀번호가 없는 경우 (소셜 로그인 유저)
		// 일반 로그인을 시도했으므로 자격 증명이 틀린 것으로 처리합니다.
		return nil, errors.New("user exists but has no password")
	}

	// 👇 password가 유효하면(NULL이 아니면) 실제 string 값을 user 구조체에 할당합니다.
	user.HashedPassword = hashedPassword.String
	return &user, nil
}


// FindOrCreateUserByGithub는 github ID로 사용자를 찾거나 새로 생성합니다.
func (s *Store) FindOrCreateUserByGithub(githubUser *models.GitHubUserResponse) (int, error) {
	var userID int
	sqlFind := `SELECT id FROM users WHERE github_id = $1`
	err := s.DB.QueryRow(sqlFind, githubUser.ID).Scan(&userID)

	if err == sql.ErrNoRows {
		// pgx 드라이버는 *string 타입의 nil 값을 DB의 NULL로 잘 처리해줍니다.
		sqlCreate := `INSERT INTO users (username, github_id, email) VALUES ($1, $2, $3) RETURNING id`
		err = s.DB.QueryRow(sqlCreate, githubUser.Login, githubUser.ID, githubUser.Email).Scan(&userID)
		if err != nil {
			return 0, err
		}
		return userID, nil
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}