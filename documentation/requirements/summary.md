Smart Company Discovery Assistant — Requirements Summary

Scenario
- Internal tool to manage company Q&A in PostgreSQL and answer natural‑language questions using an LLM.

Deliverables
- Go web service (backend + UI) with clean structure.
- Flask microservice for LLM operations (embeddings + generation).
- PostgreSQL schema and migrations for Q&A storage.
- Environment-based configuration (.env and clear README instructions).

Knowledge Base (Q&A) Management
- Table: qa_pairs(id SERIAL PRIMARY KEY, question TEXT, answer TEXT).
- Create, Edit, Delete Q&A entries with validation (non-empty).
- Delete requires confirmation.
- List view with pagination, sorting, and search.

Flask LLM Service
- Endpoint: POST /generate-answer with body {"question": "..."}.
- Use vector embeddings to find top 3 relevant Q&A pairs from PostgreSQL.
- Construct a prompt using retrieved context and generate an answer with an LLM.
- Return JSON: {"answer": "...", "status": 200}.

Go → Flask Integration
- Endpoint: POST /api/ask.
- Forwards the question to Flask, parses response, and returns as JSON or renders in UI.
- Handle network/timeouts gracefully.

UI Requirements
- Navigation bar/side bar to move between pages.
- Ask page (/ask): input, submit, and area to display answer. Calls /api/ask.
- Q&A management: forms for CRUD; list with pagination, sorting, search; validation.

Configuration & Setup
- Env vars for: database connection, Flask base URL, LLM API keys.
- Provide .env.example.
- Instructions to init DB, run both services, and seed data.

Bonus
- Docker Compose to run Go, Flask, and PostgreSQL together.

Submission
- GitHub repo with clear README; reviewers must be able to run the project from the README alone.
