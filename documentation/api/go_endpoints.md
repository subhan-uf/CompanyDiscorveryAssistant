Go Service Endpoints

UI Pages
- GET /              — Home with navigation
- GET /ask           — Ask question page
- GET /qa            — Q&A list with pagination, sorting, search (query params: page, page_size, sort, search)
- GET /qa/create     — New Q&A form
- POST /qa/create    — Create Q&A
- GET /qa/edit       — Edit form (query: id)
- POST /qa/edit      — Update Q&A (query: id)
- GET /qa/delete     — Delete confirmation (query: id)
- POST /qa/delete    — Perform delete (query: id)

APIs
- POST /api/ask      — {"question": "..."} → forwards to Flask /generate-answer and returns {"answer": "..."}
