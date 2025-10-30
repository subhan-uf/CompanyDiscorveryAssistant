package models

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type QAPair struct {
    ID       int64
    Question string
    Answer   string
}

type QAModel struct {
    DB *pgxpool.Pool
}

var ErrNotFound = errors.New("qa pair not found")

func (m *QAModel) Create(ctx context.Context, q, a string) (int64, error) {
    q = strings.TrimSpace(q)
    a = strings.TrimSpace(a)
    if q == "" || a == "" {
        return 0, fmt.Errorf("question and answer required")
    }
    var id int64
    err := m.DB.QueryRow(ctx, `INSERT INTO qa_pairs(question, answer) VALUES($1,$2) RETURNING id`, q, a).Scan(&id)
    return id, err
}

func (m *QAModel) Update(ctx context.Context, id int64, q, a string) error {
    q = strings.TrimSpace(q)
    a = strings.TrimSpace(a)
    if q == "" || a == "" {
        return fmt.Errorf("question and answer required")
    }
    ct, err := m.DB.Exec(ctx, `UPDATE qa_pairs SET question=$1, answer=$2 WHERE id=$3`, q, a, id)
    if err != nil {
        return err
    }
    if ct.RowsAffected() == 0 {
        return ErrNotFound
    }
    return nil
}

func (m *QAModel) Delete(ctx context.Context, id int64) error {
    ct, err := m.DB.Exec(ctx, `DELETE FROM qa_pairs WHERE id=$1`, id)
    if err != nil {
        return err
    }
    if ct.RowsAffected() == 0 {
        return ErrNotFound
    }
    return nil
}

type ListParams struct {
    Page     int
    PageSize int
    Sort     string // question_asc | question_desc | id_desc | id_asc
    Search   string
}

type ListResult struct {
    Items      []QAPair
    TotalCount int
}

func (m *QAModel) List(ctx context.Context, p ListParams) (ListResult, error) {
    if p.Page <= 0 {
        p.Page = 1
    }
    if p.PageSize <= 0 || p.PageSize > 100 {
        p.PageSize = 10
    }
    order := "id DESC"
    switch p.Sort {
    case "question_asc":
        order = "question ASC, id DESC"
    case "question_desc":
        order = "question DESC, id DESC"
    case "id_asc":
        order = "id ASC"
    case "id_desc":
        order = "id DESC"
    }
    where := ""
    args := []any{}
    if s := strings.TrimSpace(p.Search); s != "" {
        where = "WHERE question ILIKE $1 OR answer ILIKE $1"
        args = append(args, "%"+s+"%")
    }
    // total
    var total int
    countSQL := "SELECT COUNT(*) FROM qa_pairs " + where
    if err := m.DB.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
        return ListResult{}, err
    }
    // rows
    offset := (p.Page - 1) * p.PageSize
    listSQL := fmt.Sprintf("SELECT id, question, answer FROM qa_pairs %s ORDER BY %s LIMIT %d OFFSET %d", where, order, p.PageSize, offset)
    rows, err := m.DB.Query(ctx, listSQL, args...)
    if err != nil {
        return ListResult{}, err
    }
    defer rows.Close()
    res := ListResult{TotalCount: total}
    for rows.Next() {
        var item QAPair
        if err := rows.Scan(&item.ID, &item.Question, &item.Answer); err != nil {
            return ListResult{}, err
        }
        res.Items = append(res.Items, item)
    }
    return res, nil
}

func (m *QAModel) Get(ctx context.Context, id int64) (QAPair, error) {
    var item QAPair
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    err := m.DB.QueryRow(ctx, `SELECT id, question, answer FROM qa_pairs WHERE id=$1`, id).Scan(&item.ID, &item.Question, &item.Answer)
    if err != nil {
        return QAPair{}, err
    }
    return item, nil
}
