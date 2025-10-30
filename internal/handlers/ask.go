package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

type AskRoutes struct {
    Renderer *TemplateRenderer
    FlaskURL string
}

func (h *AskRoutes) AskPage(w http.ResponseWriter, r *http.Request) {
    h.Renderer.Render(w, "layout", map[string]any{
        "ContentTemplate": "content_ask",
    })
}

type askReq struct {
    Question string `json:"question"`
}
type askRes struct {
    Answer string `json:"answer"`
    Status int    `json:"status"`
}

func (h *AskRoutes) AskAPI(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var req askReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }
    body, _ := json.Marshal(req)
    ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
    defer cancel()
    url := h.FlaskURL + "/generate-answer"
    httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
    httpReq.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        http.Error(w, `{"error":"flask service unavailable"}`, http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()
    b, _ := io.ReadAll(resp.Body)
    if resp.StatusCode >= 300 {
        http.Error(w, string(b), http.StatusBadGateway)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(b)
}
