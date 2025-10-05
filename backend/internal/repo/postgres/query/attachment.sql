-- name: GetAttachments :many
SELECT *
FROM attachment
WHERE trip_id = $1;

-- name: GetAttachmentByID :one
SELECT *
FROM attachment
WHERE trip_id = $1
  AND id = $2;

-- name: InsertAttachment :one
INSERT INTO attachment (trip_id, name, blob)
VALUES ($1, $2, $3)
RETURNING id;

-- name: DeleteAttachmentByID :one
DELETE
FROM attachment
WHERE trip_id = $1
  AND id = $2
RETURNING id;
