CREATE TABLE IF NOT EXISTS entries (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type        TEXT    NOT NULL,
    title       TEXT    NOT NULL DEFAULT '',
    description TEXT    NOT NULL DEFAULT '',
    category    TEXT    NOT NULL DEFAULT '',
    valence     TEXT,
    recorded_at TEXT    NOT NULL,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX IF NOT EXISTS idx_entries_type      ON entries(type);
CREATE INDEX IF NOT EXISTS idx_entries_recorded  ON entries(recorded_at);
