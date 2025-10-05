-- name: InsertLocation :one
INSERT INTO location (latitude, longitude)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateLocation :exec
UPDATE location
SET latitude  = $2,
    longitude = $3
WHERE id = $1;

-- name: GetLocationIDByActivityID :one
SELECT l.id
FROM location l
         JOIN activity a on a.location_id = l.id
WHERE a.id = $1;

-- name: GetLocationIDByAccommodationID :one
SELECT l.id
FROM location l
         JOIN accommodation a on a.location_id = l.id
WHERE a.id = $1;

-- name: DeleteLocation :exec
DELETE FROM location
WHERE id = $1;

