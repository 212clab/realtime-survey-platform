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