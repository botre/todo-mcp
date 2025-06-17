CREATE TABLE IF NOT EXISTS todo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- name: GetAllPendingTodos :many
SELECT id, title, completed, created_at
FROM todo
WHERE completed = 0
ORDER BY created_at ASC;

-- name: GetAllCompletedTodos :many
SELECT id, title, completed, created_at
FROM todo
WHERE completed = 1
ORDER BY created_at DESC;

-- name: CreateTodo :one
INSERT INTO todo (title)
VALUES (?)
RETURNING id, title, completed, created_at;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = ?;

-- name: GetTodoById :one
SELECT id, title, completed, created_at
FROM todo
WHERE id = ?;

-- name: CompleteTodo :one
UPDATE todo
SET completed = 1
WHERE id = ?
RETURNING id, title, completed, created_at;

-- name: DeleteAllTodos :exec
DELETE FROM todo;
