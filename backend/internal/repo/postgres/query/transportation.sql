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

-- TODO: check tripID!
-- name: UpdateTransportation :exec
UPDATE transportation
SET type           = $2,
    origin_id      = $3,
    destination_id = $4,
    departure_time = $5,
    arrival_time   = $6,
    price          = $7
WHERE id = $1;

-- name: UpsertGenericTransportationDetail :exec
INSERT INTO transportation_generic (transportation_id, name, origin_address, destination_address)
VALUES ($1, $2, $3, $4)
ON CONFLICT(transportation_id) DO UPDATE SET
    name = $2,
    origin_address = $3,
    destination_address = $4;

-- name: GetGenericDetailByTransportationID :one
SELECT *
FROM transportation_generic
WHERE transportation_id = $1;

-- name: DeleteTransportationByID :one
DELETE
FROM transportation
WHERE trip_id = $1
  AND id = $2
RETURNING id;

-- name: GetAllGeoJson :many
SELECT transportation_geojson.geojson
FROM transportation_geojson
         JOIN transportation on transportation_geojson.transportation_id = transportation.id
WHERE transportation.trip_id = $1;

-- name: UpsertGeoJson :exec
INSERT INTO transportation_geojson (transportation_id, geojson)
VALUES ($1, $2)
ON CONFLICT(transportation_id) DO UPDATE SET
    geojson = $2;
