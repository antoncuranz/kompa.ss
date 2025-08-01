-- name: InsertLocation :one
INSERT INTO location (
    latitude, longitude
) VALUES (
    $1, $2
 )
RETURNING id;
