-- name: GetUsers :many
SELECT *
FROM "user";

-- name: GetUserByID :one
SELECT *
FROM "user"
WHERE id = $1;

-- name: GetUserByJwtSub :one
SELECT *
FROM "user"
WHERE jwt_sub = $1;

-- name: InsertUser :one
INSERT INTO "user" (name, jwt_sub)
VALUES ($1, $2)
RETURNING id;
