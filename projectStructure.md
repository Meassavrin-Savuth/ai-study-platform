# 🧠 AI Study Intelligence Platform

> A full-stack AI-powered study platform that transforms YouTube lectures and PDF documents into structured learning materials using LLMs — enabling summaries, flashcards, quizzes, and interactive Q&A for efficient learning.

---

## 📌 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Build Phases](#build-phases)
- [API Reference](#api-reference)
- [Environment Variables](#environment-variables)
- [Getting Started](#getting-started)
- [Roadmap](#roadmap)

---

## Overview

Students consume huge amounts of content from YouTube lectures and PDF textbooks — but processing it efficiently is hard. This platform solves that by turning raw content into **structured learning intelligence**.

**Core flow:**
```
YouTube Link / PDF Upload
        ↓
Content Extraction (transcript / text)
        ↓
Go Backend Processing
        ↓
AI Engine (Gemini / OpenAI)
        ↓
Structured Learning Output
        ↓
Next.js Dashboard
```

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 14 (App Router), TypeScript, Tailwind CSS |
| Backend | Go (REST API) |
| AI | Google Gemini API (gemini-1.5-flash) |
| Database | PostgreSQL (Phase 4+) |
| Deployment | Vercel (frontend) + Railway (backend) |

---

## Project Structure

```
ai-study-platform/
│
├── frontend/                        # Next.js App
│   ├── app/
│   │   ├── page.tsx                 # Home — YouTube input UI
│   │   ├── layout.tsx               # Root layout
│   │   ├── dashboard/
│   │   │   └── page.tsx             # Saved materials (Phase 4)
│   │   └── globals.css
│   │
│   ├── components/
│   │   ├── UrlInput.tsx             # YouTube URL input form
│   │   ├── SummaryCard.tsx          # Display AI summary
│   │   ├── FlashcardDeck.tsx        # Flashcard UI (Phase 6)
│   │   ├── QuizEngine.tsx           # Quiz UI (Phase 6)
│   │   └── ChatBox.tsx              # RAG chat interface (Phase 5)
│   │
│   ├── lib/
│   │   └── api.ts                   # All fetch calls to Go backend
│   │
│   ├── .env.local                   # Frontend env vars
│   └── package.json
│
├── backend/                         # Go REST API
│   ├── main.go                      # Entry point, route registration
│   │
│   ├── handlers/
│   │   ├── youtube.go               # POST /api/summarize (YouTube)
│   │   ├── pdf.go                   # POST /api/pdf (Phase 2)
│   │   └── chat.go                  # POST /api/chat (Phase 5)
│   │
│   ├── services/
│   │   ├── ai.go                    # Gemini API calls
│   │   ├── transcript.go            # YouTube transcript extraction
│   │   ├── pdf.go                   # PDF text extraction (Phase 2)
│   │   └── embeddings.go            # Vector embeddings (Phase 5)
│   │
│   ├── models/
│   │   ├── user.go                  # User model (Phase 4)
│   │   └── summary.go               # Summary model (Phase 4)
│   │
│   ├── db/
│   │   └── postgres.go              # DB connection (Phase 4)
│   │
│   ├── go.mod
│   ├── go.sum
│   └── .env                         # Backend env vars
│
└── README.md
```

---

## Build Phases

### 🟢 Phase 0 — Foundation Setup *(done)*
- [x] Next.js frontend running
- [x] Go backend running
- [x] `/health` endpoint working
- [x] Basic project structure

### 🟢 Phase 1 — YouTube → Summary (MVP)
- [ ] YouTube URL input on frontend
- [ ] Go handler receives URL
- [ ] Extract YouTube transcript
- [ ] Send transcript to Gemini API
- [ ] Return and display summary

**Key files:** `handlers/youtube.go`, `services/transcript.go`, `services/ai.go`

---

### 🟡 Phase 2 — PDF Upload
- [ ] File upload UI (frontend)
- [ ] Go handler receives PDF
- [ ] Extract text from PDF
- [ ] Send to Gemini → get summary + key points

**Key files:** `handlers/pdf.go`, `services/pdf.go`

---

### 🟡 Phase 3 — Structured Learning Output
- [ ] Update AI prompt to return JSON
- [ ] Parse structured response
- [ ] Display: summary, key points, important terms, simple explanation

**Example AI response:**
```json
{
  "summary": "...",
  "key_points": ["...", "..."],
  "important_terms": [{"term": "...", "definition": "..."}],
  "simple_explanation": "..."
}
```

---

### 🔵 Phase 4 — Database + Save System
- [ ] PostgreSQL setup
- [ ] Save summaries per user
- [ ] Google login (optional)
- [ ] Personal dashboard

**Key files:** `db/postgres.go`, `models/user.go`, `models/summary.go`

---

### 🔵 Phase 5 — Chat with Content (RAG)
- [ ] Chunk transcript/PDF text
- [ ] Generate embeddings (pgvector)
- [ ] Retrieve relevant chunks on user question
- [ ] AI answers using context

**Key files:** `services/embeddings.go`, `handlers/chat.go`

---

### 🔴 Phase 6 — Flashcards + Quiz Engine
- [ ] AI generates flashcards (JSON)
- [ ] AI generates quiz questions (JSON)
- [ ] Answer checking logic
- [ ] Score display

---

### 🔴 Phase 7 — Exam Mode
- [ ] Timed quiz mode
- [ ] Difficulty levels (easy / medium / hard)
- [ ] Progress tracking
- [ ] Score history

---

## API Reference

| Method | Endpoint | Description | Phase |
|--------|----------|-------------|-------|
| GET | `/health` | Health check | 0 |
| POST | `/api/summarize` | YouTube URL → summary | 1 |
| POST | `/api/pdf` | PDF upload → summary | 2 |
| POST | `/api/chat` | Question → RAG answer | 5 |
| GET | `/api/summaries` | Get saved summaries | 4 |

### POST `/api/summarize`
```json
// Request
{ "url": "https://youtube.com/watch?v=..." }

// Response
{
  "summary": "...",
  "key_points": ["...", "..."],
  "important_terms": [{"term": "...", "definition": "..."}]
}
```

---

## Environment Variables

### `backend/.env`
```
GEMINI_API_KEY=your_gemini_api_key
PORT=8080
DATABASE_URL=postgresql://user:password@localhost:5432/aistudy
```

### `frontend/.env.local`
```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

---

## Getting Started

### Prerequisites
- Node.js v18+
- Go v1.21+
- PostgreSQL (Phase 4+)

### Run locally

```bash
# Clone
git clone https://github.com/yourname/ai-study-platform
cd ai-study-platform

# Backend
cd backend
cp .env.example .env       # add your API keys
go run main.go             # runs on :8080

# Frontend (new terminal)
cd frontend
npm install
cp .env.local.example .env.local
npm run dev                # runs on :3000
```

### Test the API
```bash
curl -X POST http://localhost:8080/api/summarize \
  -H "Content-Type: application/json" \
  -d '{"url": "https://youtube.com/watch?v=dQw4w9WgXcQ"}'
```

---

## Roadmap

| Phase | Feature | Status |
|-------|---------|--------|
| 0 | Project setup | ✅ Done |
| 1 | YouTube → Summary | 🔨 In progress |
| 2 | PDF Upload | ⏳ Upcoming |
| 3 | Structured output | ⏳ Upcoming |
| 4 | Database + Auth | ⏳ Upcoming |
| 5 | RAG Chat | ⏳ Upcoming |
| 6 | Flashcards + Quiz | ⏳ Upcoming |
| 7 | Exam Mode | ⏳ Upcoming |

---

## Why This Project

This platform simulates real-world AI products used in EdTech, enterprise knowledge systems, and AI learning assistants. Building it covers:

- ✔ Full-stack engineering (Next.js + Go)
- ✔ Real AI integration (not just a chatbot wrapper)
- ✔ System design thinking
- ✔ Scalable architecture patterns
- ✔ Practical real-world use case

---

*Built as a portfolio project to demonstrate AI engineering, backend development, and product thinking.*