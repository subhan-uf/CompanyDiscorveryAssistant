Database Schema

Table: qa_pairs

CREATE TABLE IF NOT EXISTS qa_pairs (
  id SERIAL PRIMARY KEY,
  question TEXT NOT NULL,
  answer   TEXT NOT NULL
);

Notes
- Validation ensures non-empty question and answer at the app layer.
- Indexes can be added later if needed for search performance.
