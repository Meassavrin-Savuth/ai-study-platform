package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/yourname/ai-study-backend/services"
)

type SummarizeRequest struct {
	URL string `json:"url"`
}

type SummarizeResponse struct {
	Summary string `json:"summary"`
	Error   string `json:"error,omitempty"`
}

type transcriptResult struct {
	Transcript string `json:"transcript"`
	Error      string `json:"error"`
}

func SummarizeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Method not allowed"})
		return
	}

	var req SummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Invalid request"})
		return
	}

	if req.URL == "" {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "URL is required"})
		return
	}

	// Step 1 — get transcript
	cmd := exec.Command("python3", "scripts/get_transcript.py", req.URL)
	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Failed to get transcript"})
		return
	}

	var result transcriptResult
	if err := json.Unmarshal(output, &result); err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Failed to parse transcript"})
		return
	}

	if result.Error != "" {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: result.Error})
		return
	}

	// Step 2 — send to Gemini
	summary, err := services.Summarize(result.Transcript)
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "AI summarization failed"})
		return
	}

	json.NewEncoder(w).Encode(SummarizeResponse{Summary: summary})
}