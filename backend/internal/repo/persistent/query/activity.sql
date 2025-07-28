-- name: GetActivityByID :one
SELECT *
FROM activity
WHERE id = $1;

-- name: GetActivities :many
SELECT *
FROM activity;
