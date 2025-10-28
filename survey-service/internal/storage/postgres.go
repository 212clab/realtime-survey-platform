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

	// 👇 이 부분을 수정하세요!
	// var surveys []models.Survey // nil 슬라이스를 만듭니다.
	surveys := make([]models.Survey, 0) // 비어있는, nil이 아닌 슬라이스를 만듭니다.
	// 👆 이렇게 하면 결과가 없어도 JSON으로 `[]`가 됩니다.

	for rows.Next() {
		var survey models.Survey
		var optionsJSON []byte

		if err := rows.Scan(&survey.ID, &survey.Title, &optionsJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(optionsJSON, &survey.Options); err != nil {
			return nil, err
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}