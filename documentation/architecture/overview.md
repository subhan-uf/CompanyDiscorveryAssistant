Architecture Overview

Components
- Go Web App: HTTP server, HTML templates for UI, REST API for /api/ask and Q&A management. Connects to PostgreSQL.
- Flask LLM Service: Embeddings + answer generation via hosted providers (Cohere for embeddings, Groq/Together for chat); fetches Q&A from PostgreSQL and returns an answer.
- PostgreSQL: Stores Q&A pairs (question, answer).

Data Flow
1) CRUD: Go app ↔ PostgreSQL for managing qa_pairs.
2) Ask flow: Go app receives user question → forwards to Flask /generate-answer → Flask queries PostgreSQL, computes embeddings, selects top 3 Q&A, generates answer → Flask returns JSON → Go returns/render.

Configuration
- Env vars control DB connection string, Flask base URL, and hosted LLM API keys (Groq, Together, Cohere; OpenAI optional).

Decisions
- UI implemented with Go html/template for simplicity and minimal dependencies.
- Hosted-first LLM approach: Cohere embeddings by default (or Together), chat via Groq (or Together). OpenAI supported as a secondary option. No local model requirement.

Error Handling
- Timeouts and errors from Flask call produce user-friendly messages; server logs keep details.
