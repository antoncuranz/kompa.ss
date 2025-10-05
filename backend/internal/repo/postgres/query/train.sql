-- name: GetTrainLegsByTransportationID :many
SELECT sqlc.embed(train_leg),
       sqlc.embed(origin),
       sqlc.embed(destination),
       sqlc.embed(origin_location),
       sqlc.embed(destination_location)
FROM train_leg
         JOIN train_station origin on train_leg.origin = origin.id
         JOIN train_station destination on train_leg.destination = destination.id
         JOIN location origin_location on origin.location_id = origin_location.id
         JOIN location destination_location on destination.location_id = destination_location.id
WHERE transportation_id = $1
ORDER BY departure_time;

-- name: GetTrainDetailByTransportationID :one
SELECT *
FROM train_detail
WHERE transportation_id = $1;

-- name: InsertTrainDetail :exec
INSERT INTO train_detail (transportation_id, refresh_token)
VALUES ($1, $2);

-- name: TrainStationExists :one
SELECT EXISTS (
    SELECT 1
    FROM train_station
    WHERE id = $1
);

-- name: InsertTrainStation :exec
INSERT INTO train_station (id, name, location_id)
VALUES ($1, $2, $3);

-- name: InsertTrainLeg :one
INSERT INTO train_leg (transportation_id, origin, destination, departure_time, arrival_time, duration_in_minutes, line_name, operator_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;
