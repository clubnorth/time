package main

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "path/filepath"

  _ "modernc.org/sqlite"

  "time-server/handler"
  "time-server/middleware"
  "time-server/service"
)

type Config struct {
  Port   int    `json:"port"`
  DBPath string `json:"db_path"`
}

func loadConfig() Config {
  cfg := Config{Port: 8080, DBPath: "./data/time.db"}
  data, err := os.ReadFile("config.json")
  if err != nil {
    log.Println("config.json not found, using defaults")
    return cfg
  }
  if err := json.Unmarshal(data, &cfg); err != nil {
    log.Println("config.json parse error, using defaults")
    return cfg
  }
  return cfg
}

func main() {
  cfg := loadConfig()

  // Ensure data directory exists
  dbDir := filepath.Dir(cfg.DBPath)
  if err := os.MkdirAll(dbDir, 0755); err != nil {
    log.Fatalf("Failed to create data directory: %v", err)
  }

  // Open database
  db, err := sql.Open("sqlite", cfg.DBPath)
  if err != nil {
    log.Fatalf("Failed to open database: %v", err)
  }
  defer db.Close()

  // Run schema
  schema, err := os.ReadFile("db/schema.sql")
  if err != nil {
    log.Fatalf("Failed to read schema: %v", err)
  }
  if _, err := db.Exec(string(schema)); err != nil {
    log.Fatalf("Failed to run schema: %v", err)
  }

  // Initialize services and handlers
  entrySvc := service.NewEntryService(db)
  entryH := handler.NewEntryHandler(entrySvc)

  // Setup router
  mux := http.NewServeMux()
  mux.HandleFunc("GET /api/entries", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("month") != "" {
			entryH.GetEntries(w, r)
		} else {
			entryH.GetAllEntries(w, r)
		}
	})
  mux.HandleFunc("POST /api/entries", entryH.CreateEntry)

  // Apply middleware
  h := middleware.Logger(middleware.CORS(mux))

  addr := fmt.Sprintf(":%d", cfg.Port)
  log.Printf("Server starting on http://localhost%s", addr)
  if err := http.ListenAndServe(addr, h); err != nil {
    log.Fatalf("Server failed: %v", err)
  }
}
