-- name: GetFlightLegsByTransportationID :many
SELECT sqlc.embed(flight_leg),
       sqlc.embed(origin),
       sqlc.embed(destination),
       sqlc.embed(origin_location),
       sqlc.embed(destination_location)
FROM flight_leg
         JOIN airport origin on flight_leg.origin = origin.iata
         JOIN airport destination on flight_leg.destination = destination.iata
         JOIN location origin_location on origin.location_id = origin_location.id
         JOIN location destination_location on destination.location_id = destination_location.id
WHERE transportation_id = $1
ORDER BY departure_time;

-- name: GetPnrsByTransportationID :many
SELECT *
FROM flight_pnr
WHERE transportation_id = $1;

-- name: InsertFlightLeg :one
INSERT INTO flight_leg (transportation_id, origin, destination, airline, flight_number, departure_time, arrival_time, duration_in_minutes, aircraft)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: UpdateFlightLeg :exec
UPDATE flight_leg
SET origin              = $2,
    destination         = $3,
    airline             = $4,
    flight_number       = $5,
    departure_time      = $6,
    arrival_time        = $7,
    duration_in_minutes = $8,
    aircraft            = $9
WHERE id = $1;

-- name: AirportExists :one
SELECT EXISTS (
    SELECT 1
    FROM airport
    WHERE iata = $1
);

-- name: InsertAirport :exec
INSERT INTO airport (iata, name, municipality, location_id)
VALUES ($1, $2, $3, $4);

-- name: InsertPNR :one
INSERT INTO flight_pnr (transportation_id, airline, pnr)
VALUES ($1, $2, $3)
RETURNING id;
