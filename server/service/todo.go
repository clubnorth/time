package service

import "time"

type Todo struct {
  ID        int    `json:"id"`
  Title     string `json:"title"`
  Completed bool   `json:"completed"`
  Category  string `json:"category"`
  DueDate   string `json:"due_date"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

func (s *EntryService) GetTodos() ([]Todo, error) {
  rows, err := s.db.Query("SELECT id, title, completed, category, due_date, created_at, updated_at FROM todos ORDER BY created_at DESC")
  if err != nil { return nil, err }
  defer rows.Close()
  var todos []Todo
  for rows.Next() {
    var t Todo; var comp int
    if err := rows.Scan(&t.ID, &t.Title, &comp, &t.Category, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil { return nil, err }
    t.Completed = comp != 0
    todos = append(todos, t)
  }
  if todos == nil { todos = []Todo{} }
  return todos, nil
}

func (s *EntryService) CreateTodo(title string, category string, dueDate string) (*Todo, error) {
  now := time.Now().Format("2006-01-02 15:04:05")
  if category == "" { category = "today" }
  res, err := s.db.Exec(
    "INSERT INTO todos (title, category, due_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
    title, category, dueDate, now, now,
  )
  if err != nil { return nil, err }
  id, _ := res.LastInsertId()
  return &Todo{ID: int(id), Title: title, Completed: false, Category: category, DueDate: dueDate, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *EntryService) UpdateTodo(id int, title *string, completed *bool, category *string, dueDate *string) error {
  if title == nil && completed == nil && category == nil && dueDate == nil {
    return nil
  }
  query := "UPDATE todos SET "
  args := make([]interface{}, 0)
  if title != nil {
    query += "title = ?, "
    args = append(args, *title)
  }
  if completed != nil {
    comp := 0
    if *completed { comp = 1 }
    query += "completed = ?, "
    args = append(args, comp)
  }
  if category != nil {
    query += "category = ?, "
    args = append(args, *category)
  }
  if dueDate != nil {
    query += "due_date = ?, "
    args = append(args, *dueDate)
  }
  query += "updated_at = datetime('now','localtime') WHERE id = ?"
  args = append(args, id)
  _, err := s.db.Exec(query, args...)
  return err
}

func (s *EntryService) DeleteTodo(id int) error {
  _, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id)
  return err
}
