-- name: CreateUser :one
INSERT INTO users (id, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email;

-- name: GetUserByID :one
SELECT id, email
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email
FROM users
WHERE email = $1;