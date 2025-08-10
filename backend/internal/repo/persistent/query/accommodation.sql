-- name: GetAllAccommodation :many
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
         LEFT JOIN location on location_id = location.id
WHERE trip_id = $1;

-- name: GetAccommodationByID :one
SELECT sqlc.embed(accommodation), location.*
FROM accommodation
         LEFT JOIN location on location_id = location.id
WHERE trip_id = $1
  AND accommodation.id = $2;

-- name: InsertAccommodation :one
INSERT INTO accommodation (trip_id, location_id, name, arrival_date, departure_date, check_in_time, check_out_time, description, address, price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: UpdateAccommodation :exec
UPDATE accommodation
SET location_id    = $2,
    name           = $3,
    arrival_date   = $4,
    departure_date = $5,
    check_in_time  = $6,
    check_out_time = $7,
    description    = $8,
    address        = $9,
    price          = $10
WHERE id = $1;

-- name: DeleteAccommodationByID :exec
DELETE
FROM accommodation
WHERE trip_id = $1
  AND accommodation.id = $2;
