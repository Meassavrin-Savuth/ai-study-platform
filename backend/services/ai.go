package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type StudyMaterial struct {
	Summary           string       `json:"summary"`
	KeyPoints         []string     `json:"key_points"`
	ImportantTerms    []Term       `json:"important_terms"`
	SimpleExplanation string       `json:"simple_explanation"`
}

type Term struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

// cleanJSON strips markdown code fences that Gemini sometimes wraps around JSON
func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		// remove first line (```json or ```)
		raw = raw[strings.Index(raw, "\n")+1:]
	}
	if strings.HasSuffix(raw, "```") {
		raw = raw[:strings.LastIndex(raw, "```")]
	}
	return strings.TrimSpace(raw)
}

func callGemini(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
		apiKey,
	)

	body := GeminiRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: prompt}}},
		},
	}

	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Gemini response:", string(respBody))

	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini. Raw response: %s", string(respBody))
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func Summarize(transcript string) (*StudyMaterial, error) {
	prompt := fmt.Sprintf(`You are a study assistant. Analyze this content and return ONLY a JSON object with no markdown, no backticks, no explanation.

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
%s`, transcript)

	raw, err := callGemini(prompt)
	if err != nil {
		return nil, err
	}

	var material StudyMaterial
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &material); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %s", raw)
	}

	return &material, nil
}

func SummarizeWithPrompt(prompt string) (*StudyMaterial, error) {
	raw, err := callGemini(prompt)
	if err != nil {
		return nil, err
	}

	var material StudyMaterial
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &material); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %s", raw)
	}

	return &material, nil
}
type Flashcard struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func GenerateFlashcards(material *StudyMaterial) ([]Flashcard, error) {
	prompt := fmt.Sprintf(`You are a study assistant. Based on this study material, generate 5 flashcards.
Return ONLY a JSON array with no markdown, no backticks, no explanation.

Return exactly this structure:
[
  {"question": "Question here?", "answer": "Answer here"},
  {"question": "Question here?", "answer": "Answer here"},
  {"question": "Question here?", "answer": "Answer here"},
  {"question": "Question here?", "answer": "Answer here"},
  {"question": "Question here?", "answer": "Answer here"}
]

Study material:
Summary: %s
Key points: %v
Important terms: %v`, material.Summary, material.KeyPoints, material.ImportantTerms)

	raw, err := callGemini(prompt)
	if err != nil {
		return nil, err
	}

	var flashcards []Flashcard
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &flashcards); err != nil {
		fmt.Printf("Flashcard parse error: %v\nRaw: %s\n", err, raw)
		return nil, fmt.Errorf("failed to parse flashcards: %s", raw)
	}

	return flashcards, nil
}