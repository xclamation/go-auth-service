-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, ip_address)
VALUES ($1, $2, $3, $4)
RETURNING id, email;

-- name: GetUserByID :one
SELECT id, email, ip_address
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email
FROM users
WHERE email = $1;