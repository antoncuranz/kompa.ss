-- name: GetTrips :many
SELECT trip.*
FROM trip
LEFT JOIN permissions p ON trip.id = p.trip_id
WHERE trip.owner_id = sqlc.arg(user_id) OR p.user_id = sqlc.arg(user_id);

-- name: GetTripByID :one
SELECT *
FROM trip
WHERE id = $1;

-- name: InsertTrip :one
INSERT INTO trip (owner_id, name, start_date, end_date, description, image_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateTrip :exec
UPDATE trip
SET name        = $2,
    start_date  = $3,
    end_date    = $4,
    description = $5,
    image_url   = $6
WHERE id = $1;

-- name: DeleteTripByID :one
DELETE
FROM trip
WHERE id = $1
RETURNING id;
