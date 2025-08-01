-- name: GetActivityByID :one
SELECT sqlc.embed(activity), location.*
FROM activity
LEFT JOIN location on activity.location_id = location.id
WHERE activity.id = $1;

-- name: GetActivities :many
SELECT sqlc.embed(activity), location.*
FROM activity
LEFT JOIN location on activity.location_id = location.id;
