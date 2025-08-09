-- name: GetActivityByID :one
SELECT sqlc.embed(activity), location.*
FROM activity
LEFT JOIN location on activity.location_id = location.id
WHERE trip_id = $1
AND activity.id = $2;

-- name: GetActivities :many
SELECT sqlc.embed(activity), location.*
FROM activity
LEFT JOIN location on activity.location_id = location.id
WHERE trip_id = $1;

-- name: InsertActivity :one
INSERT INTO activity (
    trip_id, location_id, name, date, time, address, description, price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;
