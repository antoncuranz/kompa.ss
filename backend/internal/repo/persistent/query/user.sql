-- name: GetUserByID :one
SELECT *
FROM "user"
WHERE id = $1;

-- name: GetUsers :many
SELECT *
FROM "user";
