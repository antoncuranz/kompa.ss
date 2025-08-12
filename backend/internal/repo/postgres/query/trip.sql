-- name: GetTrips :many
SELECT *
FROM trip;

-- name: GetTripByID :one
SELECT *
FROM trip
WHERE id = $1;

-- name: InsertTrip :one
INSERT INTO trip (name, start_date, end_date, description, image_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateTrip :exec
UPDATE trip
SET name        = $2,
    start_date  = $3,
    end_date    = $4,
    description = $5,
    image_url   = $6
WHERE id = $1;

-- name: DeleteTripByID :exec
DELETE
FROM trip
WHERE id = $1;
