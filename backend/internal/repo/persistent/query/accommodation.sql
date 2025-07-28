-- name: GetAccommodationByID :one
SELECT *
FROM accommodation
WHERE id = $1;

-- name: GetAllAccommodation :many
SELECT *
FROM accommodation;
