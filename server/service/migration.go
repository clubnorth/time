package service

import (
  "database/sql"
  "fmt"
  "log"
  "os"
  "path/filepath"
  "sort"
  "strconv"
  "strings"
)

type migration struct {
  version int
  name    string
  sql     string
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
  // Ensure schema_migrations table exists
  _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name    TEXT NOT NULL,
    applied_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
  )`)
  if err != nil {
    return fmt.Errorf("failed to create schema_migrations table: %w", err)
  }

  // Get applied migrations
  applied := make(map[int]bool)
  rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
  if err != nil {
    return fmt.Errorf("failed to query applied migrations: %w", err)
  }
  defer rows.Close()
  for rows.Next() {
    var v int
    if err := rows.Scan(&v); err != nil {
      return err
    }
    applied[v] = true
  }

  // Read migration files
  entries, err := os.ReadDir(migrationsDir)
  if err != nil {
    return fmt.Errorf("failed to read migrations directory '%s': %w", migrationsDir, err)
  }

  var migrations []migration
  for _, entry := range entries {
    if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
      continue
    }
    verStr := strings.Split(entry.Name(), "_")[0]
    ver, err := strconv.Atoi(verStr)
    if err != nil {
      continue
    }
    content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
    if err != nil {
      return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
    }
    migrations = append(migrations, migration{version: ver, name: entry.Name(), sql: string(content)})
  }

  sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })

  // Run pending migrations
  for _, m := range migrations {
    if applied[m.version] {
      continue
    }
    log.Printf("Running migration %03d: %s", m.version, m.name)
    if _, err := db.Exec(m.sql); err != nil {
      return fmt.Errorf("migration %03d failed: %w", m.version, err)
    }
    if _, err := db.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", m.version, m.name); err != nil {
      return fmt.Errorf("failed to record migration %03d: %w", m.version, err)
    }
    log.Printf("Migration %03d completed", m.version)
  }

  return nil
}

// ExportData dumps all data as JSON-serializable structures
type ExportData struct {
  Entries  []Entry        `json:"entries"`
  Todos    []Todo         `json:"todos"`
  Settings []ExportSetting `json:"settings"`
}

type ExportSetting struct {
  Key   string `json:"key"`
  Value string `json:"value"`
}

func (s *EntryService) ExportAll() (*ExportData, error) {
  data := &ExportData{}

  entries, err := s.GetEntriesPaginated(100000, "")
  if err != nil {
    return nil, err
  }
  if entries == nil {
    entries = []Entry{}
  }
  data.Entries = entries

  todos, err := s.GetTodos()
  if err != nil {
    return nil, err
  }
  if todos == nil {
    todos = []Todo{}
  }
  data.Todos = todos

  rows, err := s.db.Query("SELECT key, value FROM settings")
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  for rows.Next() {
    var s ExportSetting
    if err := rows.Scan(&s.Key, &s.Value); err != nil {
      return nil, err
    }
    data.Settings = append(data.Settings, s)
  }
  if data.Settings == nil {
    data.Settings = []ExportSetting{}
  }

  return data, nil
}

func (s *EntryService) ImportAll(data *ExportData) error {
  // Clear existing data
  tx, err := s.db.Begin()
  if err != nil {
    return err
  }
  defer tx.Rollback()

  if _, err := tx.Exec("DELETE FROM entries"); err != nil {
    return err
  }
  if _, err := tx.Exec("DELETE FROM todos"); err != nil {
    return err
  }
  if _, err := tx.Exec("DELETE FROM settings"); err != nil {
    return err
  }

  // Insert entries with original IDs and timestamps
  for _, e := range data.Entries {
    _, err := tx.Exec(
      "INSERT INTO entries (id, type, title, description, category, valence, recorded_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
      e.ID, e.Type, e.Title, e.Description, e.Category, e.Valence, e.RecordedAt, e.CreatedAt, e.UpdatedAt,
    )
    if err != nil {
      return fmt.Errorf("failed to import entry %d: %w", e.ID, err)
    }
  }

  // Insert todos with original IDs and timestamps
  for _, t := range data.Todos {
    comp := 0
    if t.Completed {
      comp = 1
    }
    _, err := tx.Exec(
      "INSERT INTO todos (id, title, completed, category, due_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
      t.ID, t.Title, comp, t.Category, t.DueDate, t.CreatedAt, t.UpdatedAt,
    )
    if err != nil {
      return fmt.Errorf("failed to import todo %d: %w", t.ID, err)
    }
  }

  // Insert settings
  for _, s := range data.Settings {
    _, err := tx.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", s.Key, s.Value)
    if err != nil {
      return fmt.Errorf("failed to import setting %s: %w", s.Key, err)
    }
  }

  return tx.Commit()
}
