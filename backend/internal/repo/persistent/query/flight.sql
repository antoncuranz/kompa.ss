-- name: GetFlightByID :one
SELECT *
FROM flight f
WHERE trip_id = $1
AND id = $2;

-- name: GetFlights :many
SELECT *
FROM flight
WHERE trip_id = $1;

-- name: GetFlightLegsByFlightID :many
SELECT sqlc.embed(flight_leg), sqlc.embed(origin), sqlc.embed(destination), sqlc.embed(origin_location), sqlc.embed(destination_location)
FROM flight_leg
JOIN airport origin on flight_leg.origin = origin.iata
JOIN airport destination on flight_leg.destination = destination.iata
JOIN location origin_location on origin.location_id = origin_location.id
JOIN location destination_location on destination.location_id = destination_location.id
WHERE flight_id = $1
ORDER BY departure_time;

-- name: GetPnrsByFlightID :many
SELECT *
FROM pnr
WHERE flight_id = $1;

-- name: InsertFlight :one
INSERT INTO flight (
    trip_id, price
) VALUES (
    $1, $2
 )
RETURNING id;

-- name: InsertFlightLeg :one
INSERT INTO flight_leg (
    flight_id, origin, destination, airline, flight_number, departure_time, arrival_time, duration_in_minutes, aircraft
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id;

-- name: InsertAirport :exec
INSERT INTO airport (
    iata, name, municipality, location_id
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT DO NOTHING;

-- name: InsertPNR :one
INSERT INTO pnr (
    flight_id, airline, pnr
) VALUES (
    $1, $2, $3
)
RETURNING id;
