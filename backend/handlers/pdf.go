package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/yourname/ai-study-backend/services"
)

type pdfResult struct {
	Text  string `json:"text"`
	Error string `json:"error"`
}

func PDFHandler(w http.ResponseWriter, r *http.Request) {
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

	// Max 10MB file
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("pdf")
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "No PDF file received"})
		return
	}
	defer file.Close()

	// Save to temp file
	tmp, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Failed to save file"})
		return
	}
	defer os.Remove(tmp.Name())

	io.Copy(tmp, file)
	tmp.Close()

	// Extract text with Python
	cmd := exec.Command("python3", "scripts/get_pdf_text.py", tmp.Name())
	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Failed to extract PDF text"})
		return
	}

	var result pdfResult
	if err := json.Unmarshal(output, &result); err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "Failed to parse PDF text"})
		return
	}

	if result.Error != "" {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: result.Error})
		return
	}

	// Limit text size to avoid token limits
	text := result.Text
	if len(text) > 8000 {
		text = text[:8000]
	}

	// Send to Gemini
	prompt := fmt.Sprintf(`You are a study assistant. Summarize this PDF document clearly for a student.

Include:
- A 3-5 sentence summary
- 3-5 key points as bullet points
- 3 important terms with definitions

Document:
%s`, text)

	summary, err := services.SummarizeWithPrompt(prompt)
	if err != nil {
		json.NewEncoder(w).Encode(SummarizeResponse{Error: "AI summarization failed"})
		return
	}

	json.NewEncoder(w).Encode(SummarizeResponse{Summary: summary})
}