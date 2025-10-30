Requirements Compliance

Core System Setup
- Go service: Implemented under `cmd/server` with modular packages in `internal/*`.
- Flask microservice: Implemented in `flask_service/app.py`.
- PostgreSQL schema/migration: `migrations/001_init.sql` auto-applied on Go server startup.
- Env/config: `.env.example` provided; README documents setup.

Knowledge Base Management (Go)
- Table `qa_pairs(id SERIAL PRIMARY KEY, question TEXT, answer TEXT)` created in migration.
- UI pages:
  - Create: `GET /qa/create`, `POST /qa/create` with validation (non-empty fields)
  - Edit: `GET /qa/edit?id=...`, `POST /qa/edit?id=...` with validation
  - Delete: `GET /qa/delete?id=...` shows confirmation, `POST /qa/delete?id=...` performs delete
  - List: `GET /qa` supports pagination (`page`, `page_size`), sorting (`sort`), and search (`search`)

Flask LLM Service
- Endpoint: `POST /generate-answer` with `{ "question": "..." }` implemented.
- Embeddings retrieval: Uses vector embeddings to select top 3 relevant pairs.
  - Primary: Cohere embeddings (hosted). Alternative: Together embeddings.
  - Secondary: OpenAI embeddings (optional).
  - Fallback: Simple token overlap (only if embedding backend is unavailable).
- Answer generation: Builds prompt from top 3 Q&A and generates with:
  - Primary: Groq Llama 3.1 (hosted). Alternative: Together chat.
  - Secondary: OpenAI chat (optional).
  - Fallback: Returns the best-matching stored answer.

Go → Flask Integration
- Go endpoint `POST /api/ask` forwards to Flask `/generate-answer`, returns response or meaningful error on failure.

UI Requirements
- Navigation between Home, Ask, and Q&A.
- Ask page `/ask`: input box, submit, and answer display; calls `/api/ask`.
- Q&A page `/qa`: CRUD and list with pagination, sorting, search.

Configuration & Setup
- Env documented in `.env.example` and `documentation/setup/env.md`.
- README provides DB init, service run, and seed instructions.

Bonus (Optional)
- Docker Compose provided (`docker-compose.yml`) to run Postgres, Flask, and Go together.
