package handler

import (
  "encoding/json"
  "net/http"
  "strconv"
)

func (h *EntryHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
  todos, err := h.svc.GetTodos()
  if err != nil { respondError(w, 500, err.Error()); return }
  respond(w, 200, todos)
}

func (h *EntryHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
  var body struct {
    Title    string `json:"title"`
    Category string `json:"category"`
    DueDate  string `json:"due_date"`
  }
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
    respondError(w, 400, "title required"); return
  }
  todo, err := h.svc.CreateTodo(body.Title, body.Category, body.DueDate)
  if err != nil { respondError(w, 500, err.Error()); return }
  respond(w, 201, todo)
}

func (h *EntryHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil { respondError(w, 400, "invalid id"); return }
  var body struct {
    Title     *string `json:"title"`
    Completed *bool   `json:"completed"`
    Category  *string `json:"category"`
    DueDate   *string `json:"due_date"`
  }
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    respondError(w, 400, "invalid body"); return
  }
  if err := h.svc.UpdateTodo(id, body.Title, body.Completed, body.Category, body.DueDate); err != nil {
    respondError(w, 500, err.Error()); return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}

func (h *EntryHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil { respondError(w, 400, "invalid id"); return }
  if err := h.svc.DeleteTodo(id); err != nil {
    respondError(w, 500, err.Error()); return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}
