-- name: GetTrips :many
SELECT trip.*
FROM trip
LEFT JOIN permissions p ON trip.id = p.trip_id
WHERE trip.owner_id = sqlc.arg(user_id) OR p.user_id = sqlc.arg(user_id);

-- name: GetTripByID :one
SELECT trip.*
FROM trip
LEFT JOIN permissions p ON trip.id = p.trip_id
WHERE (trip.owner_id = sqlc.arg(user_id) OR p.user_id = sqlc.arg(user_id))
AND id = $1;

-- name: InsertTrip :one
INSERT INTO trip (owner_id, name, start_date, end_date, description, image_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateTrip :exec
UPDATE trip
SET name        = $3,
    start_date  = $4,
    end_date    = $5,
    description = $6,
    image_url   = $7
WHERE trip.owner_id = $2
AND id = $1;

-- name: DeleteTripByID :exec
DELETE
FROM trip
WHERE trip.owner_id = $2
AND id = $1;
