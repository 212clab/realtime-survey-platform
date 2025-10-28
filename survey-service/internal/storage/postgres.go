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

// GetAllSurveys는 모든 설문조사 목록을 데이터베이스에서 가져옵니다.
func (s *Store) GetAllSurveys() ([]models.Survey, error) {
	rows, err := s.DB.Query("SELECT id, title, options FROM surveys ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []models.Survey
	for rows.Next() {
		var survey models.Survey
		var optionsJSON []byte // options를 JSON byte 슬라이스로 받음

		if err := rows.Scan(&survey.ID, &survey.Title, &optionsJSON); err != nil {
			return nil, err
		}

		// JSON byte 슬라이스를 models.Option 슬라이스로 변환
		if err := json.Unmarshal(optionsJSON, &survey.Options); err != nil {
			return nil, err
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}