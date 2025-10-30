-- Create qa_pairs table
CREATE TABLE IF NOT EXISTS qa_pairs (
  id SERIAL PRIMARY KEY,
  question TEXT NOT NULL,
  answer   TEXT NOT NULL
);
