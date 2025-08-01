-- name: GetAccommodationByID :one
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
LEFT JOIN location on location_id = location.id
WHERE accommodation.id = $1;

-- name: GetAllAccommodation :many
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
LEFT JOIN location on location_id = location.id;
