package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"

	"github.com/yourname/ai-study-backend/handlers"
)

func main() {
	godotenv.Load() 

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	http.HandleFunc("/api/summarize", handlers.SummarizeHandler)
	http.HandleFunc("/api/pdf", handlers.PDFHandler)
	http.HandleFunc("/api/flashcards", handlers.FlashcardsHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}