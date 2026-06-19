package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourname/ai-study-backend/services"
)

type FlashcardRequest struct {
	Summary        string          `json:"summary"`
	KeyPoints      []string        `json:"key_points"`
	ImportantTerms []services.Term `json:"important_terms"`
}

func FlashcardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req FlashcardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	material := &services.StudyMaterial{
		Summary:        req.Summary,
		KeyPoints:      req.KeyPoints,
		ImportantTerms: req.ImportantTerms,
	}

	flashcards, err := services.GenerateFlashcards(material)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate flashcards"})
		return
	}

	json.NewEncoder(w).Encode(flashcards)
}