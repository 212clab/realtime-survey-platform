package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"user-service/internal/api"      // api 패키지 import
	"user-service/internal/storage"  // storage 패키지 import

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// 1. DB 연결
	dbSource := os.Getenv("DB_SOURCE")
	db, err := sql.Open("pgx", dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("✅ Successfully connected to the database!")

	// 2. 테이블 생성 (DB 초기화)
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// 3. 의존성 조립 (storage -> api)
	store := storage.NewStore(db)
	api := api.NewAPI(store)

	// 4. 라우터 설정
	http.HandleFunc("/signup", api.SignupHandler)
	http.HandleFunc("/health", api.HealthCheckHandler)

	// 5. 서버 시작
	fmt.Println("✅ User service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}