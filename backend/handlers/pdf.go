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
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("pdf")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "No PDF file received"})
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}
	defer os.Remove(tmp.Name())

	io.Copy(tmp, file)
	tmp.Close()

	cmd := exec.Command("python3", "scripts/get_pdf_text.py", tmp.Name())
	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to extract PDF text"})
		return
	}

	var result pdfResult
	if err := json.Unmarshal(output, &result); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse PDF text"})
		return
	}

	if result.Error != "" {
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error})
		return
	}

	text := result.Text
	if len(text) > 8000 {
		text = text[:8000]
	}

	prompt := fmt.Sprintf(`You are a study assistant. Analyze this PDF and return ONLY a JSON object with no markdown, no backticks, no explanation.

Return exactly this structure:
{
  "summary": "3-5 sentence summary here",
  "key_points": ["point 1", "point 2", "point 3", "point 4", "point 5"],
  "important_terms": [
    {"term": "Term 1", "definition": "Definition here"},
    {"term": "Term 2", "definition": "Definition here"},
    {"term": "Term 3", "definition": "Definition here"}
  ],
  "simple_explanation": "Explain this like I am 15 years old in 2-3 sentences"
}

Content:
%s`, text)

	material, err := services.SummarizeWithPrompt(prompt)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "AI summarization failed"})
		return
	}

	json.NewEncoder(w).Encode(material)
}