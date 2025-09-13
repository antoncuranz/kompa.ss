-- name: GetAllTransportation :many
SELECT sqlc.embed(transportation),
       sqlc.embed(origin),
       sqlc.embed(destination)
FROM transportation
         JOIN location origin on transportation.origin_id = origin.id
         JOIN location destination on transportation.destination_id = destination.id
WHERE transportation.trip_id = $1;

-- name: GetTransportationByID :one
SELECT sqlc.embed(transportation),
       sqlc.embed(origin),
       sqlc.embed(destination)
FROM transportation
         JOIN location origin on transportation.origin_id = origin.id
         JOIN location destination on transportation.destination_id = destination.id
WHERE transportation.trip_id = $1
  AND transportation.id = $2;

-- name: InsertTransportation :one
INSERT INTO transportation (trip_id, type, origin_id, destination_id, departure_time, arrival_time, geojson, price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: DeleteTransportationByID :exec
DELETE
FROM transportation
WHERE trip_id = $1
  AND transportation.id = $2;
