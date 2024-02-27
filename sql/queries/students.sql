-- name: CreateStudent :one
INSERT INTO students(id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetStudents :many
SELECT * FROM students;