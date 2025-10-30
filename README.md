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
1) Install dependencies
   - Go 1.22+
   - Python 3.10+
   - PostgreSQL 14+ (local optional; we use Supabase hosted Postgres in these steps)

2) Configure database (Supabase hosted)
   - Copy env example:
     cp .env.example .env
   - Set `DATABASE_URL` to your Supabase Postgres URI with sslmode=require
     e.g., postgres://postgres:<password>@db.<project-ref>.supabase.co:5432/postgres?sslmode=require

3) Initialize schema and seed (optional)
   - The Go server applies `migrations/001_init.sql` automatically at startup.
   - To seed sample Q&A in Supabase: open SQL Editor → paste `scripts/seed.sql` → Run

4) Start the Flask LLM service
   - Create a virtual environment and install deps:
     cd flask_service && python3 -m venv .venv && source .venv/bin/activate && pip install -r requirements.txt
   - Hosted providers (no local models):
     - Sign up and get free API keys:
       Groq (chat): https://console.groq.com/keys
       Together (chat/embeddings): https://api.together.xyz
       Cohere (embeddings): https://dashboard.cohere.com/api-keys
     - Export keys (any combination works):
       export GROQ_API_KEY=grq_...
       export TOGETHER_API_KEY=...
       export COHERE_API_KEY=...
   - Run Flask:
     export DATABASE_URL=${DATABASE_URL}
     flask --app app.py run --host 0.0.0.0 --port 5000

5) Start the Go web app
   - In a new terminal at repo root:
     export $(grep -v '^#' .env | xargs) || true
     go run ./cmd/server
   - Open http://localhost:8080

Usage
- Navigate to Q&A to create/edit/delete entries. The list supports search, sorting, and pagination via controls at the top.
- Open Ask to submit a question. The Go API calls the Flask service, which retrieves relevant Q&A and generates an answer.

Configuration
- DATABASE_URL: PostgreSQL connection (e.g., postgres://user:pass@localhost:5432/smartassistant?sslmode=disable)
- FLASK_URL: Flask service base URL (default http://localhost:5000)
- PORT: Go server port (default 8080)
- GROQ_API_KEY (Flask optional): Hosted Llama chat via Groq
- TOGETHER_API_KEY (Flask optional): Hosted Llama chat and embeddings via Together
- COHERE_API_KEY (Flask optional): Hosted embeddings via Cohere
- OPENAI_API_KEY (Flask optional): Secondary backend (not required)

Testing The Flow
- Seed data via scripts/seed.sql
- Start both services as described
- Ask a question like: "What is the refund policy?" and verify the answer.

Migrations
- Startup applies migrations/001_init.sql. For further changes, add new sequential SQL files.

Troubleshooting
- Flask service unavailable: ensure Flask is running and FLASK_URL is correct.
- DB connection errors: verify DATABASE_URL and that Postgres is up.
- OpenAI disabled: If OPENAI_API_KEY not set, retrieval falls back to simple matching.

Documentation
- See documentation/ for architecture, API contracts, schema, and env details.
- Requirements compliance: documentation/requirements/compliance.md
 - Supabase setup (hosted Postgres): documentation/setup/supabase.md

Docker (Optional)
- Start everything with Docker Compose:
  - docker compose up --build
  - Open http://localhost:8080 for the Go app, Flask at http://localhost:5000
  - To enable OpenAI in Docker, export OPENAI_API_KEY before running compose.
