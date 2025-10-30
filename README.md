Smart Company Discovery Assistant

Overview
- Multi-service app to manage internal Q&A and answer user questions using LLMs.
- Components: Go web app (UI + API), Flask LLM service, PostgreSQL database.

Requirements Recap
- CRUD for Q&A entries with validation and delete confirmation
- Ask page that calls /api/ask
- Flask service uses embeddings to retrieve top‑3 relevant Q&A and generates an answer
- Clear setup using environment variables

Quick Start (Local)
This project assumes you already have a valid `.env` file provided to you (DB + API keys). No additional configuration is required.

1) Install dependencies
   - Go 1.22+
   - Python 3.10+

2) Start the Flask LLM service
   - cd flask_service && python3 -m venv .venv && source .venv/bin/activate && pip install -r requirements.txt
   - export $(grep -v '^#' ../.env | xargs)
   - flask --app app.py run --host 0.0.0.0 --port 5000

3) Start the Go web app
   - In a new terminal at repo root:
     export $(grep -v '^#' .env | xargs) || true
     go run ./cmd/server
   - Open http://localhost:8080

4) Initialize schema and seed (optional)
   - On first start, the Go server applies `migrations/001_init.sql` automatically.
   - To seed sample Q&A: in Supabase SQL Editor, paste `scripts/seed.sql` and Run

Usage
- Navigate to Q&A to create/edit/delete entries. The list supports search, sorting, and pagination via controls at the top.
- Open Ask to submit a question. The Go API calls the Flask service, which retrieves relevant Q&A and generates an answer.

Configuration
- All required environment variables are provided via `.env`.
- Common variables: `DATABASE_URL`, `FLASK_URL`, `PORT`, and LLM provider keys.

Testing The Flow
- Seed data via scripts/seed.sql
- Start both services as described
- Ask a question like: "What is the refund policy?" and verify the answer.

Migrations
- Startup applies `migrations/001_init.sql`. For further changes, add new sequential SQL files.

Troubleshooting
- Flask unavailable: ensure Flask is running and `FLASK_URL` (from `.env`) is correct.
- DB connection: verify `DATABASE_URL` in `.env`.

Documentation
- See documentation/ for architecture, API contracts, schema, and env details.
- Requirements compliance: documentation/requirements/compliance.md
 - Supabase setup (hosted Postgres): documentation/setup/supabase.md

Docker (Optional)
- docker compose up --build
- Go app at http://localhost:8080, Flask at http://localhost:5000
