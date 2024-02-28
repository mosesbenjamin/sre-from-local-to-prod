-- name: CreateStudent :one
INSERT INTO students(id, created_at, updated_at, email, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetStudents :many
SELECT * FROM students ORDER BY created_at DESC LIMIT $1;

-- name: GetStudentById :one
SELECT * FROM students WHERE id = $1;

-- name: GetStudentByEmail :one
SELECT * FROM students WHERE email = $1;

-- name: DeleteStudent :exec
DELETE FROM students WHERE id = $1;

-- name: UpdateStudentPassword :one
UPDATE students SET password = $2 WHERE id = $1 RETURNING *;
