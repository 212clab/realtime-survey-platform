package api

import (
	"encoding/json"
	"net/http"
	"survey-service/internal/models"
	"survey-service/internal/storage"
)

type API struct {
	Store *storage.Store
}

func NewAPI(store *storage.Store) *API {
	return &API{Store: store}
}

// SurveysHandler는 /surveys 경로의 요청을 HTTP 메서드에 따라 분기합니다.
func (a *API) SurveysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.GetAllSurveysHandler(w, r)
	case http.MethodPost:
		a.CreateSurveyHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllSurveysHandler는 모든 설문조사 목록을 반환합니다.
func (a *API) GetAllSurveysHandler(w http.ResponseWriter, r *http.Request) {
	surveys, err := a.Store.GetAllSurveys()
	if err != nil {
		http.Error(w, "Failed to get surveys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(surveys)
}



// CreateSurveyHandler는 설문 생성 요청을 처리합니다.
func (a *API) CreateSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var survey models.Survey
	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	surveyID, err := a.Store.CreateSurvey(&survey)
	if err != nil {
		http.Error(w, "Failed to create survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Survey created successfully",
		"surveyId": surveyID,
	})
}

func (a *API) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}