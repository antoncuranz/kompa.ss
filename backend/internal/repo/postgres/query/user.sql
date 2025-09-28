-- name: GetUsers :many
SELECT *
FROM "user";

-- name: GetUserByID :one
SELECT *
FROM "user"
WHERE id = $1;

-- name: GetUserByJwtSub :one
SELECT *
FROM "user"
WHERE jwt_sub = $1;

-- name: InsertUser :one
INSERT INTO "user" (name, jwt_sub)
VALUES ($1, $2)
RETURNING id;

-- name: HasReadPermission :one
SELECT EXISTS (
    SELECT 1
    FROM trip
             LEFT JOIN permissions p ON trip.id = p.trip_id
    WHERE (trip.owner_id = sqlc.arg(user_id) OR p.user_id = sqlc.arg(user_id))
      AND id = sqlc.arg(trip_id)
);

-- name: HasWritePermission :one
SELECT EXISTS (
    SELECT 1
    FROM trip
             LEFT JOIN permissions p ON trip.id = p.trip_id
    WHERE (trip.owner_id = sqlc.arg(user_id) OR (p.user_id = sqlc.arg(user_id) AND p.write is true))
      AND id = sqlc.arg(trip_id)
);

-- name: IsTripOwner :one
SELECT EXISTS (
    SELECT 1
    FROM trip
    WHERE trip.owner_id = sqlc.arg(user_id)
      AND id = sqlc.arg(trip_id)
);
