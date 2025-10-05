-- name: GetActivities :many
SELECT sqlc.embed(activity), location.*
FROM activity
         LEFT JOIN location on activity.location_id = location.id
WHERE trip_id = $1;

-- name: GetActivityByID :one
SELECT sqlc.embed(activity), location.*
FROM activity
         LEFT JOIN location on activity.location_id = location.id
WHERE trip_id = $1
  AND activity.id = $2;

-- name: InsertActivity :one
INSERT INTO activity (trip_id, location_id, name, date, time, address, description, price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;

-- TODO: check tripID!
-- name: UpdateActivity :exec
UPDATE activity
SET location_id = $2,
    name        = $3,
    date        = $4,
    time        = $5,
    address     = $6,
    description = $7,
    price       = $8
WHERE id = $1;

-- name: DeleteActivityByID :exec
DELETE
FROM activity
WHERE trip_id = $1
  AND activity.id = $2;
