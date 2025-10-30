Smart Company Discovery Assistant

Overview
- Multi‑service app for managing internal Q&A and answering questions using an LLM.
- Components: Go web app (UI + API), Flask LLM service, PostgreSQL database.

Quick Start (Local)
This project assumes you already have a valid `.env` provided by the maintainer. No extra key or database setup is required.

1) Install prerequisites
- Go 1.22+
- Python 3.10+

2) Start services (choose your OS)

Linux / macOS
- Flask service
  - `cd flask_service`
  - `python3 -m venv .venv && source .venv/bin/activate`
  - `pip install -r requirements.txt`
  - `export $(grep -v '^#' ../.env | xargs)`
  - `flask --app app.py run --host 0.0.0.0 --port 5000`
- Go web app (new terminal at repo root)
  - `export $(grep -v '^#' .env | xargs) || true`
  - `go run ./cmd/server`
  - Open http://localhost:8080

Windows (PowerShell)
- Flask service
  - `cd flask_service`
  - `py -3 -m venv .venv`
  - `.\.venv\Scripts\Activate.ps1`
  - `pip install -r requirements.txt`
  - Load .env into the current session:
    - `Get-Content ..\.env | Where-Object { $_ -and $_ -notmatch '^\s*#' } | ForEach-Object { $p = $_ -split '=',2; if ($p.Length -eq 2) { Set-Item -Path Env:$($p[0]) -Value $p[1] } }`
  - `flask --app app.py run --host 0.0.0.0 --port 5000`
- Go web app (new PowerShell at repo root)
  - `Get-Content .\.env | Where-Object { $_ -and $_ -notmatch '^\s*#' } | ForEach-Object { $p = $_ -split '=',2; if ($p.Length -eq 2) { Set-Item -Path Env:$($p[0]) -Value $p[1] } }`
  - `go run .\cmd\server`
  - Open http://localhost:8080

3) Optional: seed sample data
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
