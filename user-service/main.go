package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL 드라이버 import
	"golang.org/x/crypto/bcrypt" // bcrypt 라이브러리 import
)

// DB 인스턴스를 전역적으로 접근할 수 있도록 구조체 선언
type application struct {
	db *sql.DB
}

// User 모델 정의
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func main() {
	// 1. 환경 변수에서 DB 연결 정보 가져오기
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable is not set")
	}

	// 2. 데이터베이스에 연결
	db, err := sql.Open("pgx", dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close() // main 함수가 끝나면 DB 연결 종료

	// 3. 실제 DB 연결을 확인 (Ping)
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	fmt.Println("✅ Successfully connected to the database!")

	// //////////////// users 테이블 생성(start)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

    if _, err := db.Exec(createTableSQL); err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }

	// app 인스턴스 생성
	app := &application{
		db: db,
	}

	// 라우터 설정
	http.HandleFunc("/signup", app.signupHandler) // 새로운 핸들러 등록
	http.HandleFunc("/health", app.healthCheckHandler)

	fmt.Println("✅ Service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// signupHandler: 회원가입 요청을 처리하는 함수
func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	// 1. 요청 본문(JSON)을 User 구조체로 디코딩
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. 비밀번호를 bcrypt로 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 3. DB에 새로운 사용자 정보 삽입
	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = app.db.Exec(sqlStatement, user.Username, string(hashedPassword))
	if err != nil {
		// 여기서는 간단히 처리하지만, 실제로는 username 중복 에러 등을 구분해야 함
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// 4. 성공 응답 전송
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// healthCheckHandler: DB 상태를 확인하는 함수 (기존 로직과 동일)
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.db.Ping(); err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "OK")
}
