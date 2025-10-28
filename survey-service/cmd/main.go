package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"survey-service/internal/api"
	"survey-service/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbSource := os.Getenv("DB_SOURCE")
	db, err := sql.Open("pgx", dbSource)
	if err != nil { log.Fatalf("DB 연결 실패: %v", err) }
	defer db.Close()
	if err := db.Ping(); err != nil { log.Fatalf("DB Ping 실패: %v", err) }
	fmt.Println("✅ DB 연결 성공!")

	// surveys 테이블 생성
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS surveys (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        options JSONB NOT NULL
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("surveys 테이블 생성 실패: %v", err)
	}

	store := storage.NewStore(db)
	api := api.NewAPI(store)

	http.HandleFunc("/surveys", api.SurveysHandler)
	http.HandleFunc("/health", api.HealthCheckHandler)

	fmt.Println("✅ Survey service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("서버 시작 실패: %s\n", err)
	}
}