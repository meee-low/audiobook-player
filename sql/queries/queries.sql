-- name: GetPerson :one
SELECT * FROM persons
WHERE id = ? LIMIT 1;

-- name: ListPersons :many
SELECT * FROM persons
ORDER BY name;

-- name: CreatePerson :one
INSERT INTO persons (
  name
) VALUES (
  ?
)
RETURNING *;

-- name: UpdatePerson :exec
UPDATE persons
set name = ?
WHERE id = ?
RETURNING *;

-- name: DeletePerson :exec
DELETE FROM persons
WHERE id = ?;


