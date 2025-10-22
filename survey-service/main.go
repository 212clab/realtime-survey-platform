package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL 드라이버 import
)

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

	// 기본 핸들러
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the Service! DB connection is successful.")
	})

	// DB 연결 상태 확인을 위한 핸들러
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "DB connection failed", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("✅ Service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}