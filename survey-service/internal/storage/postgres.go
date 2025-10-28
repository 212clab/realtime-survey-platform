package storage

import (
	"database/sql"
	"encoding/json"
	"survey-service/internal/models"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

// CreateSurvey는 새로운 설문조사를 데이터베이스에 생성합니다.
func (s *Store) CreateSurvey(survey *models.Survey) (int, error) {
	// Options를 JSON 문자열로 변환하여 저장합니다.
	optionsJSON, err := json.Marshal(survey.Options)
	if err != nil {
		return 0, err
	}

	var surveyID int
	sqlStatement := `INSERT INTO surveys (title, options) VALUES ($1, $2) RETURNING id`
	err = s.DB.QueryRow(sqlStatement, survey.Title, optionsJSON).Scan(&surveyID)
	if err != nil {
		return 0, err
	}

	return surveyID, nil
}