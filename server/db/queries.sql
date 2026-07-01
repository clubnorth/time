-- name: GetEntriesByMonth :many
SELECT * FROM entries
WHERE strftime('%Y-%m', recorded_at) = ?
ORDER BY recorded_at DESC, id DESC;

-- name: CreateEntry :one
INSERT INTO entries (type, title, description, category, valence, recorded_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;
