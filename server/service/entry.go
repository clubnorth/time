package service

import (
  "database/sql"
  "time"
)

type Entry struct {
  ID          int     `json:"id"`
  Type        string  `json:"type"`
  Title       string  `json:"title"`
  Description string  `json:"description"`
  Category    string  `json:"category"`
  Valence     *string `json:"valence"`
  RecordedAt  string  `json:"recorded_at"`
  CreatedAt   string  `json:"created_at"`
  UpdatedAt   string  `json:"updated_at"`
}

type EntryService struct {
  db *sql.DB
}

func NewEntryService(db *sql.DB) *EntryService {
  return &EntryService{db: db}
}

func (s *EntryService) GetEntriesByMonth(month string) ([]Entry, error) {
  query := `SELECT id, type, title, description, category, valence, recorded_at, created_at, updated_at FROM entries WHERE strftime('%Y-%m', recorded_at) = ? ORDER BY recorded_at DESC, id DESC`
  rows, err := s.db.Query(query, month)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var entries []Entry
  for rows.Next() {
    var e Entry
    err := rows.Scan(&e.ID, &e.Type, &e.Title, &e.Description, &e.Category, &e.Valence, &e.RecordedAt, &e.CreatedAt, &e.UpdatedAt)
    if err != nil {
      return nil, err
    }
    entries = append(entries, e)
  }
  return entries, nil
}


func (s *EntryService) GetEntriesPaginated(limit int, before string) ([]Entry, error) {
  var query string
  var args []interface{}

  if before == "" {
    query = "SELECT id, type, title, description, category, valence, recorded_at, created_at, updated_at FROM entries ORDER BY recorded_at DESC, id DESC LIMIT ?"
    args = []interface{}{limit}
  } else {
    query = "SELECT id, type, title, description, category, valence, recorded_at, created_at, updated_at FROM entries WHERE recorded_at < ? ORDER BY recorded_at DESC, id DESC LIMIT ?"
    args = []interface{}{before, limit}
  }

  rows, err := s.db.Query(query, args...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var entries []Entry
  for rows.Next() {
    var e Entry
    err := rows.Scan(&e.ID, &e.Type, &e.Title, &e.Description, &e.Category, &e.Valence, &e.RecordedAt, &e.CreatedAt, &e.UpdatedAt)
    if err != nil {
      return nil, err
    }
    entries = append(entries, e)
  }
  return entries, nil
}

func (s *EntryService) CreateEntry(e *Entry) error {
  now := time.Now().Format("2006-01-02 15:04:05")
  query := `INSERT INTO entries (type, title, description, category, valence, recorded_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
  result, err := s.db.Exec(query, e.Type, e.Title, e.Description, e.Category, e.Valence, e.RecordedAt, now, now)
  if err != nil {
    return err
  }
  id, err := result.LastInsertId()
  if err != nil {
    return err
  }
  e.ID = int(id)
  e.CreatedAt = now
  e.UpdatedAt = now
  return nil
}
