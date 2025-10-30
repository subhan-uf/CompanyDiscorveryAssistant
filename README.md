Smart Company Discovery Assistant

Overview
- Multi‑service app for managing internal Q&A and answering questions using an LLM.
- Components: Go web app (UI + API), Flask LLM service, PostgreSQL database.

Quick Start (Local)
This project assumes you already have a valid `.env` provided by the maintainer. No extra key or database setup is required.

1) Install prerequisites
- Go 1.22+
- Python 3.10+

2) Start the Flask LLM service
- cd flask_service
- python3 -m venv .venv && source .venv/bin/activate
- pip install -r requirements.txt
- export $(grep -v '^#' ../.env | xargs)
- flask --app app.py run --host 0.0.0.0 --port 5000

3) Start the Go web app
- In a new terminal at the repo root:
  - export $(grep -v '^#' .env | xargs) || true
  - go run ./cmd/server
- Open http://localhost:8080

4) Optional: seed sample data
- On first start, the Go server applies `migrations/001_init.sql` automatically.
- To seed examples: use your DB console (e.g., Supabase SQL Editor), paste `scripts/seed.sql`, then Run.

Usage
- Q&A Management: http://localhost:8080/qa (list, search, sort, pagination, create/edit/delete)
- Ask: http://localhost:8080/ask (type a question and view the generated answer)

Configuration
- All required environment variables are read from the provided `.env`.

Migrations
- Startup applies `migrations/001_init.sql`. For further schema changes, add a new sequential SQL file in `migrations/`.

Troubleshooting
- Flask unavailable: ensure Flask is running and `FLASK_URL` from `.env` points to http://localhost:5000 (default).
- DB connection: verify `DATABASE_URL` in `.env`.

Documentation
- See `documentation/` for architecture, endpoints, env, and requirements compliance.

Docker (Optional)
- docker compose up --build
- Go app: http://localhost:8080, Flask: http://localhost:5000
