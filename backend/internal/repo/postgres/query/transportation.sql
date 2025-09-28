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
INSERT INTO transportation (trip_id, type, origin_id, destination_id, departure_time, arrival_time, price)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: DeleteTransportationByID :exec
DELETE
FROM transportation
WHERE trip_id = $1
  AND transportation.id = $2;

-- name: GetAllGeoJson :many
SELECT transportation_geojson.geojson
FROM transportation_geojson
         JOIN transportation on transportation_geojson.transportation_id = transportation.id
WHERE transportation.trip_id = $1;

-- name: InsertGeoJson :exec
INSERT INTO transportation_geojson (transportation_id, geojson)
VALUES ($1, $2)
    ON CONFLICT(transportation_id)
DO UPDATE SET
    geojson = $2;
