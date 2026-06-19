package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/yourname/ai-study-backend/services"
)

type SummarizeRequest struct {
	URL string `json:"url"`
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req SummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if req.URL == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
		return
	}

	cmd := exec.Command("python3", "scripts/get_transcript.py", req.URL)
	output, err := cmd.Output()
	if err != nil {
		// capture stderr for better diagnostics
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Transcript script stderr: %s\n", string(exitErr.Stderr))
		}
		fmt.Printf("Transcript script error: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to run transcript script"})
		return
	}

	var result transcriptResult
	if err := json.Unmarshal(output, &result); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse transcript"})
		return
	}

	if result.Error != "" {
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error})
		return
	}

	material, err := services.Summarize(result.Transcript)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "AI summarization failed"})
		return
	}

	json.NewEncoder(w).Encode(material)
}