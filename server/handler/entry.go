package handler

import (
	"strconv"
  "encoding/json"
  "net/http"
  "time-server/service"
)

type EntryHandler struct {
  svc *service.EntryService
}

func NewEntryHandler(svc *service.EntryService) *EntryHandler {
  return &EntryHandler{svc: svc}
}

type apiResponse struct {
  Code    int         `json:"code"`
  Message string      `json:"message"`
  Data    interface{} `json:"data"`
}

func respond(w http.ResponseWriter, code int, data interface{}) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(code)
  json.NewEncoder(w).Encode(apiResponse{
    Code:    0,
    Message: "ok",
    Data:    data,
  })
}

func respondError(w http.ResponseWriter, code int, msg string) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(code)
  json.NewEncoder(w).Encode(apiResponse{
    Code:    1,
    Message: msg,
    Data:    nil,
  })
}

func (h *EntryHandler) GetEntries(w http.ResponseWriter, r *http.Request) {
  month := r.URL.Query().Get("month")
  if month == "" {
    respondError(w, 400, "month parameter required")
    return
  }
  entries, err := h.svc.GetEntriesByMonth(month)
  if err != nil {
    respondError(w, 500, err.Error())
    return
  }
  if entries == nil {
    entries = []service.Entry{}
  }
  respond(w, 200, entries)
}


func (h *EntryHandler) GetAllEntries(w http.ResponseWriter, r *http.Request) {
  limitStr := r.URL.Query().Get("limit")
  before := r.URL.Query().Get("before")

  limit := 30
  if limitStr != "" {
    if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 100 {
      limit = n
    }
  }

  entries, err := h.svc.GetEntriesPaginated(limit, before)
  if err != nil {
    respondError(w, 500, err.Error())
    return
  }
  if entries == nil {
    entries = []service.Entry{}
  }
  respond(w, 200, entries)
}

func (h *EntryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
  var entry service.Entry
  if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
    respondError(w, 400, "invalid request body")
    return
  }
  if entry.Type == "" || entry.RecordedAt == "" {
    respondError(w, 400, "type and recorded_at are required")
    return
  }
  if err := h.svc.CreateEntry(&entry); err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 201, entry)
}

func (h *EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil {
    respondError(w, 400, "invalid id")
    return
  }
  if err := h.svc.DeleteEntry(id); err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}
func (h *EntryHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
  key := r.PathValue("key")
  if key == "" {
    respondError(w, 400, "key required")
    return
  }
  value, err := h.svc.GetSetting(key)
  if err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, map[string]string{"key": key, "value": value})
}

func (h *EntryHandler) SetSetting(w http.ResponseWriter, r *http.Request) {
  key := r.PathValue("key")
  if key == "" {
    respondError(w, 400, "key required")
    return
  }
  var body struct {
    Value string `json:"value"`
  }
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    respondError(w, 400, "invalid body")
    return
  }
  if err := h.svc.SetSetting(key, body.Value); err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}

func (h *EntryHandler) RecalculateEntries(w http.ResponseWriter, r *http.Request) {
  entryType := r.URL.Query().Get("type")
  if entryType == "" {
    respondError(w, 400, "type required")
    return
  }
  if err := h.svc.RecalculateEntries(entryType); err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}

func (h *EntryHandler) ExportData(w http.ResponseWriter, r *http.Request) {
  data, err := h.svc.ExportAll()
  if err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, data)
}

func (h *EntryHandler) ImportData(w http.ResponseWriter, r *http.Request) {
  var data service.ExportData
  if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
    respondError(w, 400, "invalid request body")
    return
  }
  if err := h.svc.ImportAll(&data); err != nil {
    respondError(w, 500, err.Error())
    return
  }
  respond(w, 200, map[string]string{"status": "ok"})
}
