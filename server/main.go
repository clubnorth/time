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
  Host          string `json:"host"`
  Port          int    `json:"port"`
  DBPath        string `json:"db_path"`
  StaticDir     string `json:"static_dir"`
  MigrationsDir string `json:"migrations_dir"`
  DeepSeekKey   string `json:"deepseek_key"`
}

func loadConfig() Config {
  cfg := Config{Host: "0.0.0.0", Port: 8080, DBPath: "./data/time.db", StaticDir: "./dist", MigrationsDir: "./db/migrations"}
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

  // Run migrations
  if err := service.RunMigrations(db, cfg.MigrationsDir); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
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
  mux.HandleFunc("DELETE /api/entries/{id}", entryH.DeleteEntry)
  mux.HandleFunc("POST /api/entries/recalculate", entryH.RecalculateEntries)
  mux.HandleFunc("GET /api/data/export", entryH.ExportData)
  mux.HandleFunc("POST /api/data/import", entryH.ImportData)
  mux.HandleFunc("GET /api/settings/{key}", entryH.GetSetting)
  mux.HandleFunc("PUT /api/settings/{key}", entryH.SetSetting)
  mux.HandleFunc("GET /api/todos", entryH.ListTodos)
  mux.HandleFunc("POST /api/todos", entryH.CreateTodo)
  mux.HandleFunc("PUT /api/todos/{id}", entryH.UpdateTodo)
  mux.HandleFunc("DELETE /api/todos/{id}", entryH.DeleteTodo)
  mux.HandleFunc("POST /api/book-info", handler.BookInfo(cfg.DeepSeekKey))
  mux.HandleFunc("POST /api/media-info", handler.MediaInfo(cfg.DeepSeekKey))

  // Static file serving + SPA fallback
  staticDir := cfg.StaticDir
  if _, err := os.Stat(staticDir); err == nil {
    fileServer := http.FileServer(http.Dir(staticDir))
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      path := r.URL.Path
      fullPath := filepath.Join(staticDir, path)

      // If the requested path is a file or directory that exists, serve it
      info, err := os.Stat(fullPath)
      if err == nil && !info.IsDir() {
        fileServer.ServeHTTP(w, r)
        return
      }

      // Check if the path without trailing slash is a directory
      if err == nil && info.IsDir() {
        indexPath := filepath.Join(fullPath, "index.html")
        if _, err := os.Stat(indexPath); err == nil {
          fileServer.ServeHTTP(w, r)
          return
        }
      }

      // SPA fallback: serve index.html for all non-API routes
      r.URL.Path = "/"
      fileServer.ServeHTTP(w, r)
    })
    log.Printf("Serving static files from: %s", staticDir)
  } else {
    log.Printf("Static dir not found (%s), running API-only mode", staticDir)
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path == "/" {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok","message":"API server running"}`))
        return
      }
      http.NotFound(w, r)
    })
  }

  // Apply middleware
  h := middleware.Logger(middleware.CORS(mux))

  addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
  log.Printf("Server starting on http://%s", addr)
  if err := http.ListenAndServe(addr, h); err != nil {
    log.Fatalf("Server failed: %v", err)
  }
}
