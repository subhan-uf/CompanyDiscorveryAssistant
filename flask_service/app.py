import os
import json
import math
from typing import List, Tuple

from flask import Flask, request, jsonify
import psycopg
from psycopg.rows import dict_row
import requests

try:
    from openai import OpenAI
except Exception:  # pragma: no cover
    OpenAI = None  # type: ignore


app = Flask(__name__)


DB_DSN = os.getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smartassistant?sslmode=disable")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
GROQ_API_KEY = os.getenv("GROQ_API_KEY")
TOGETHER_API_KEY = os.getenv("TOGETHER_API_KEY")
COHERE_API_KEY = os.getenv("COHERE_API_KEY")

# Ollama config (free local models)
OLLAMA_BASE_URL = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")
OLLAMA_EMBED_MODEL = os.getenv("OLLAMA_EMBED_MODEL", "nomic-embed-text")
OLLAMA_CHAT_MODEL = os.getenv("OLLAMA_CHAT_MODEL", "llama3.1:8b")
OLLAMA_ENABLED = os.getenv("OLLAMA_ENABLED", "true").lower() in ("1", "true", "yes", "on")

_conn = None
_openai_client = None


def get_db():
    global _conn
    if _conn is None:
        _conn = psycopg.connect(DB_DSN)
    return _conn


def get_openai():
    global _openai_client
    if _openai_client is None and OPENAI_API_KEY and OpenAI is not None:
        _openai_client = OpenAI(api_key=OPENAI_API_KEY)
    return _openai_client


def ollama_available() -> bool:
    if not OLLAMA_ENABLED:
        return False
    try:
        r = requests.get(f"{OLLAMA_BASE_URL}/api/tags", timeout=1.5)
        return r.status_code == 200
    except Exception:
        return False


def fetch_qas() -> List[Tuple[int, str, str]]:
    con = get_db()
    with con.cursor(row_factory=dict_row) as cur:
        cur.execute("SELECT id, question, answer FROM qa_pairs ORDER BY id DESC")
        rows = cur.fetchall()
        return [(int(r["id"]), str(r["question"] or ""), str(r["answer"] or "")) for r in rows]


def cosine_sim(a: List[float], b: List[float]) -> float:
    if not a or not b or len(a) != len(b):
        return 0.0
    dot = sum(x*y for x, y in zip(a, b))
    na = math.sqrt(sum(x*x for x in a))
    nb = math.sqrt(sum(y*y for y in b))
    if na == 0.0 or nb == 0.0:
        return 0.0
    return dot / (na * nb)


def embed_ollama(texts: List[str]) -> List[List[float]]:
    embs: List[List[float]] = []
    for t in texts:
        payload = {"model": OLLAMA_EMBED_MODEL, "prompt": t}
        r = requests.post(f"{OLLAMA_BASE_URL}/api/embeddings", json=payload, timeout=30)
        r.raise_for_status()
        data = r.json()
        embs.append(data.get("embedding") or data.get("embeddings") or [])
    return embs


def embed_openai(texts: List[str]) -> List[List[float]]:
    client = get_openai()
    if client is None:
        raise RuntimeError("OpenAI client not available")
    res = client.embeddings.create(model="text-embedding-3-small", input=texts)
    return [d.embedding for d in res.data]


def embed_cohere(texts: List[str]) -> List[List[float]]:
    if not COHERE_API_KEY:
        raise RuntimeError("COHERE_API_KEY not set")
    url = "https://api.cohere.com/v1/embed"
    headers = {"Authorization": f"Bearer {COHERE_API_KEY}", "Content-Type": "application/json"}
    payload = {"texts": texts, "model": "embed-english-v3.0"}
    r = requests.post(url, headers=headers, json=payload, timeout=60)
    r.raise_for_status()
    data = r.json()
    return data.get("embeddings", [])


def embed_together(texts: List[str]) -> List[List[float]]:
    if not TOGETHER_API_KEY:
        raise RuntimeError("TOGETHER_API_KEY not set")
    url = "https://api.together.xyz/v1/embeddings"
    headers = {"Authorization": f"Bearer {TOGETHER_API_KEY}", "Content-Type": "application/json"}
    payload = {"input": texts, "model": "togethercomputer/m2-bert-80M-8k-retrieval"}
    r = requests.post(url, headers=headers, json=payload, timeout=60)
    r.raise_for_status()
    data = r.json()
    arr = []
    for item in data.get("data", []):
        arr.append(item.get("embedding", []))
    return arr


def top3_by_embedding(question: str, qas: List[Tuple[int, str, str]]):
    # Prefer Ollama embeddings; else Cohere/Together/OpenAI; else overlap fallback
    try:
        q_texts = [f"{q}\n{a}" for (_, q, a) in qas]
        if ollama_available():
            embs = embed_ollama([question] + q_texts)
        elif COHERE_API_KEY:
            embs = embed_cohere([question] + q_texts)
        elif TOGETHER_API_KEY:
            embs = embed_together([question] + q_texts)
        else:
            embs = embed_openai([question] + q_texts)
        q_emb = embs[0]
        qa_embs = embs[1:]
        scored = []
        for i, e in enumerate(qa_embs):
            sim = cosine_sim(q_emb, e)
            scored.append((sim, qas[i]))
        scored.sort(key=lambda x: x[0], reverse=True)
        return [item for (_, item) in scored[:3]]
    except Exception:
        # Fallback simple overlap scoring
        qset = set(question.lower().split())
        scored = []
        for qa in qas:
            _, q, a = qa
            aset = set((q + " " + a).lower().split())
            inter = len(qset & aset)
            scored.append((inter, qa))
        scored.sort(key=lambda x: x[0], reverse=True)
        return [item for (_, item) in scored[:3]]


def generate_answer_ollama(question: str, context_text: str) -> str:
    messages = [
        {
            "role": "system",
            "content": (
                "You are a concise, professional internal assistant. \n"
                "Answer strictly using the provided Q&A context. If unknown, say you do not know. \n"
                "Paraphrase the answer; do not copy wording verbatim. \n"
                "Keep it short: 1–3 sentences, clear and neutral."
            ),
        },
        {"role": "user", "content": f"Context:\n{context_text}\n\nUser question: {question}"},
    ]
    payload = {"model": OLLAMA_CHAT_MODEL, "messages": messages, "stream": False}
    r = requests.post(f"{OLLAMA_BASE_URL}/api/chat", json=payload, timeout=120)
    r.raise_for_status()
    data = r.json()
    # Both /api/chat and /api/generate exist; for chat, expect 'message' with 'content'
    if isinstance(data, dict):
        msg = data.get("message") or {}
        content = (msg.get("content") if isinstance(msg, dict) else None) or data.get("response")
        if content:
            return str(content)
    return ""


def generate_answer_openai(question: str, context_text: str) -> str:
    client = get_openai()
    if client is None:
        return ""
    prompt = (
        "You are a concise, professional internal assistant. Use ONLY the following Q&A context to answer. "
        "If the context does not contain the answer, say you do not know. \n"
        "Paraphrase; do not copy wording verbatim. Reply in 1–3 sentences.\n\n"
        f"Context:\n{context_text}\n\nUser question: {question}"
    )
    chat = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {"role": "system", "content": "You answer strictly based on provided context."},
            {"role": "user", "content": prompt},
        ],
        temperature=0.2,
    )
    return chat.choices[0].message.content or ""


def generate_answer_groq(question: str, context_text: str) -> str:
    if not GROQ_API_KEY:
        return ""
    url = "https://api.groq.com/openai/v1/chat/completions"
    headers = {"Authorization": f"Bearer {GROQ_API_KEY}", "Content-Type": "application/json"}
    prompt = (
        "You are a concise, professional internal assistant. Answer strictly using the Q&A context. "
        "If unknown, say you do not know. Paraphrase; do not copy wording verbatim. "
        "Respond in 1–3 sentences.\n\n"
        f"Context:\n{context_text}\n\nUser question: {question}"
    )
    payload = {
        "model": "llama-3.1-8b-instant",
        "messages": [
            {"role": "system", "content": "You answer strictly based on provided context."},
            {"role": "user", "content": prompt},
        ],
        "temperature": 0.2,
    }
    r = requests.post(url, headers=headers, json=payload, timeout=120)
    r.raise_for_status()
    data = r.json()
    return data.get("choices", [{}])[0].get("message", {}).get("content", "")


def generate_answer_together(question: str, context_text: str) -> str:
    if not TOGETHER_API_KEY:
        return ""
    url = "https://api.together.xyz/v1/chat/completions"
    headers = {"Authorization": f"Bearer {TOGETHER_API_KEY}", "Content-Type": "application/json"}
    prompt = (
        "You are a concise, professional internal assistant. Answer strictly using the Q&A context. "
        "If unknown, say you do not know. Paraphrase; do not copy wording verbatim. "
        "Respond in 1–3 sentences.\n\n"
        f"Context:\n{context_text}\n\nUser question: {question}"
    )
    payload = {
        "model": "meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo",
        "messages": [
            {"role": "system", "content": "You answer strictly based on provided context."},
            {"role": "user", "content": prompt},
        ],
        "temperature": 0.2,
    }
    r = requests.post(url, headers=headers, json=payload, timeout=120)
    r.raise_for_status()
    data = r.json()
    return data.get("choices", [{}])[0].get("message", {}).get("content", "")


def generate_answer(question: str, context_qas: List[Tuple[int, str, str]]) -> str:
    context_text = "\n\n".join([f"Q: {q}\nA: {a}" for (_, q, a) in context_qas])
    # Prefer Ollama chat; else Groq/Together/OpenAI; else fallback to best matching answer
    try:
        if ollama_available():
            out = generate_answer_ollama(question, context_text)
            if out:
                return out
    except Exception:
        pass
    try:
        out = generate_answer_groq(question, context_text)
        if out:
            return out
    except Exception:
        pass
    try:
        out = generate_answer_together(question, context_text)
        if out:
            return out
    except Exception:
        pass
    try:
        out = generate_answer_openai(question, context_text)
        if out:
            return out
    except Exception:
        pass
    return context_qas[0][2] if context_qas else "I'm not sure based on the available Q&A."


@app.post("/generate-answer")
def generate_answer_route():
    try:
        payload = request.get_json(force=True)
    except Exception:
        return jsonify({"error": "Invalid JSON"}), 400
    question = (payload or {}).get("question", "").strip()
    if not question:
        return jsonify({"error": "question is required"}), 400

    qas = fetch_qas()
    top = top3_by_embedding(question, qas)
    answer = generate_answer(question, top)
    return jsonify({"answer": answer, "status": 200})


if __name__ == "__main__":
    host = os.getenv("FLASK_HOST", "0.0.0.0")
    port = int(os.getenv("FLASK_PORT", "5000"))
    app.run(host=host, port=port)
