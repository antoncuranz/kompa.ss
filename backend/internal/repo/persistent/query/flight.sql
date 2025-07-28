-- name: GetFlightByID :one
SELECT *
FROM flight f
WHERE id = $1;

-- name: GetFlights :many
SELECT *
FROM flight;

-- name: GetFlightLegsByFlightID :many
SELECT sqlc.embed(flight_leg), sqlc.embed(origin), sqlc.embed(destination)
FROM flight_leg
JOIN airport origin on flight_leg.origin = origin.iata
JOIN airport destination on flight_leg.destination = destination.iata
WHERE flight_id = $1;

-- name: GetPnrsByFlightID :many
SELECT *
FROM pnr
WHERE flight_id = $1;
