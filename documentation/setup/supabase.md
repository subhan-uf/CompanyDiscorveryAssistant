Supabase (Hosted Postgres) Setup

Overview
- Use Supabase as the PostgreSQL backend instead of a local Postgres server.
- The app connects via a standard Postgres URI with TLS (`sslmode=require`).

Steps
1) Create a Supabase project: https://supabase.com
2) Find your Postgres connection string:
   - Supabase → Settings → Database → Connection Info → URI
   - It looks like: `postgres://postgres:PASSWORD@db.<project-ref>.supabase.co:5432/postgres?sslmode=require`
3) Export the connection string in your shell before starting services:
   - export DATABASE_URL="postgres://postgres:PASSWORD@db.<project-ref>.supabase.co:5432/postgres?sslmode=require"
4) Run the services (see README Quick Start).
   - Go server auto-applies `migrations/001_init.sql` on first run to create `qa_pairs`.
5) Seed data (optional):
   - Supabase → SQL Editor → paste the contents of `scripts/seed.sql` → Run
   - Or use the UI at `/qa` to add entries.

Notes
- Do not commit secrets to Git. Keep them in shell or a local `.env` file.
- If you change your DB password in Supabase, update `DATABASE_URL` accordingly.
