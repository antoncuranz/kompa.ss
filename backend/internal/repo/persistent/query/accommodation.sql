-- name: GetAccommodationByID :one
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
LEFT JOIN location on location_id = location.id
WHERE trip_id = $1
AND accommodation.id = $2;

-- name: GetAllAccommodation :many
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
LEFT JOIN location on location_id = location.id
WHERE accommodation.id = $1;

-- name: InsertAccommodation :one
INSERT INTO accommodation (
    trip_id, location_id, name, arrival_date, departure_date, check_in_time, check_out_time, description, address, price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;
