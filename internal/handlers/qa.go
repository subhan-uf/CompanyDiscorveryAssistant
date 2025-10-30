package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smartassistant/internal/models"
)

type QARoutes struct {
	Renderer *TemplateRenderer
	Model    *models.QAModel
}

func (h *QARoutes) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("page_size"))
	params := models.ListParams{
		Page:     page,
		PageSize: size,
		Sort:     q.Get("sort"),
		Search:   q.Get("search"),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.Model.List(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Renderer.Render(w, "layout", map[string]any{
		"ContentTemplate": "content_qa_list",
		"Items":           res.Items,
		"Total":           res.TotalCount,
		"Params":          params,
		"QueryRaw":        r.URL.RawQuery,
		"HasSearch":       strings.TrimSpace(params.Search) != "",
		"Year":            time.Now().Year(),
	})
}

func (h *QARoutes) CreateForm(w http.ResponseWriter, r *http.Request) {
	h.Renderer.Render(w, "layout", map[string]any{
		"ContentTemplate": "content_qa_form",
		"Mode":            "create",
		"Year":            time.Now().Year(),
	})
}

func (h *QARoutes) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err := h.Model.Create(ctx, question, answer)
	if err != nil {
		h.Renderer.Render(w, "layout", map[string]any{
			"ContentTemplate": "content_qa_form",
			"Mode":            "create",
			"Error":           err.Error(),
			"Question":        question,
			"Answer":          answer,
			"Year":            time.Now().Year(),
		})
		return
	}
	http.Redirect(w, r, "/qa", http.StatusSeeOther)
}

func (h *QARoutes) EditForm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	item, err := h.Model.Get(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.Renderer.Render(w, "layout", map[string]any{
		"ContentTemplate": "content_qa_form",
		"Mode":            "edit",
		"ID":              item.ID,
		"Question":        item.Question,
		"Answer":          item.Answer,
		"Year":            time.Now().Year(),
	})
}

func (h *QARoutes) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.Model.Update(ctx, id, question, answer); err != nil {
		h.Renderer.Render(w, "layout", map[string]any{
			"ContentTemplate": "content_qa_form",
			"Mode":            "edit",
			"ID":              id,
			"Error":           err.Error(),
			"Question":        question,
			"Answer":          answer,
			"Year":            time.Now().Year(),
		})
		return
	}
	http.Redirect(w, r, "/qa", http.StatusSeeOther)
}

func (h *QARoutes) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	item, err := h.Model.Get(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.Renderer.Render(w, "layout", map[string]any{
		"ContentTemplate": "content_qa_delete",
		"ID":              item.ID,
		"Question":        item.Question,
		"Year":            time.Now().Year(),
	})
}

func (h *QARoutes) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_ = h.Model.Delete(ctx, id)
	http.Redirect(w, r, "/qa", http.StatusSeeOther)
}
