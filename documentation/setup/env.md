Environment Variables

Common
- DATABASE_URL: PostgreSQL connection string (e.g., postgres://user:pass@localhost:5432/smartassistant?sslmode=disable)

Go Service
- FLASK_URL: Base URL for the Flask service (e.g., http://localhost:5000)
- PORT: Go server port (default 8080)

Flask Service
- DATABASE_URL: PostgreSQL connection string (reused)
- OPENAI_API_KEY: Optional; used as a secondary backend
- FLASK_HOST: Flask bind host (default 0.0.0.0)
- FLASK_PORT: Flask port (default 5000)
- GROQ_API_KEY: Use Groq hosted Llama for chat (OpenAI-compatible endpoint)
- TOGETHER_API_KEY: Use Together for chat and/or embeddings
- COHERE_API_KEY: Use Cohere for embeddings (embed-english-v3.0)

Docker Compose
- db: exposes 5432, credentials postgres/postgres, db smartassistant
- flask: bound to 5000, uses DATABASE_URL pointing to db container; pass OPENAI_API_KEY from host if set
- go: bound to 8080, uses FLASK_URL=http://flask:5000 and DATABASE_URL pointing to db container
