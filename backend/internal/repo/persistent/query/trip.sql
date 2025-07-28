-- name: GetTripByID :one
SELECT *
FROM trip
WHERE id = $1;

-- name: GetTrips :many
SELECT *
FROM trip;
