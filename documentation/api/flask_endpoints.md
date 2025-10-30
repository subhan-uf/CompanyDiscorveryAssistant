Flask Service Endpoints

- POST /generate-answer
  - Request JSON: {"question": "..."}
  - Process:
    1) Connect to PostgreSQL and fetch Q&A pairs
    2) Compute embeddings: question vs. each Q&A (question+answer text)
    3) Select top 3 by cosine similarity
    4) Build prompt and generate an answer
    5) Return JSON
  - Response JSON: {"answer": "...", "status": 200}
