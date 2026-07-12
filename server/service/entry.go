package service

import (
  "database/sql"
  "strconv"
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

func (s *EntryService) GetSetting(key string) (string, error) {
  var value string
  err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
  if err == sql.ErrNoRows {
    return "", nil
  }
  return value, err
}

func (s *EntryService) SetSetting(key, value string) error {
  _, err := s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
  return err
}

func (s *EntryService) DeleteEntry(id int) error {
  _, err := s.db.Exec("DELETE FROM entries WHERE id = ?", id)
  return err
}

func (s *EntryService) UpdateEntry(id int, e *Entry) error {
  now := time.Now().Format("2006-01-02 15:04:05")
  _, err := s.db.Exec(
    "UPDATE entries SET title=?, description=?, category=?, valence=?, recorded_at=?, updated_at=? WHERE id=?",
    e.Title, e.Description, e.Category, e.Valence, e.RecordedAt, now, id,
  )
  return err
}

func (s *EntryService) RecalculateEntries(entryType string) error {
  rows, err := s.db.Query(
    "SELECT id, recorded_at FROM entries WHERE type = ? ORDER BY recorded_at ASC",
    entryType,
  )
  if err != nil {
    return err
  }
  defer rows.Close()

  type row struct {
    id         int
    recordedAt string
  }
  var entries []row
  for rows.Next() {
    var r row
    if err := rows.Scan(&r.id, &r.recordedAt); err != nil {
      return err
    }
    entries = append(entries, r)
  }
  if len(entries) == 0 {
    return nil
  }

  for i := len(entries) - 1; i >= 0; i-- {
    count := 1
    for j := i - 1; j >= 0; j-- {
      d1, _ := time.Parse("2006-01-02", entries[j+1].recordedAt[:10])
      d2, _ := time.Parse("2006-01-02", entries[j].recordedAt[:10])
      if d1.Sub(d2).Hours() >= 23 && d1.Sub(d2).Hours() <= 25 {
        count++
      } else {
        break
      }
    }
    _, err := s.db.Exec("UPDATE entries SET description = ? WHERE id = ?", strconv.Itoa(count), entries[i].id)
    if err != nil {
      return err
    }
  }
  return nil
}
