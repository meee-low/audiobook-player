-- name: GetPerson :one
SELECT * FROM persons
WHERE id = ? LIMIT 1;

-- name: GetPersonByName :one
SELECT * FROM persons
WHERE name = ? LIMIT 1;


-- name: CreatePerson :one
INSERT INTO persons (
  name
) VALUES (
  ?
) RETURNING *;

-- name: CreateParentPath :one
INSERT INTO paths ( path )
VALUES ( ? )
RETURNING *;

-- name: GetPathByName :one
SELECT * FROM paths
WHERE path = ? LIMIT 1;


-- name: CreateFileInfo :one
INSERT INTO files (
  name, extension, seconds, parent_path_id, checksum, filesize
) VALUES (
  ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetFileByFullPath :one
SELECT * FROM files
WHERE (
  parent_path_id = ? AND
  name = ? AND
  extension = ?
) LIMIT 1;


-- name: CreateBook :one
INSERT INTO books (
  title, release_year
) VALUES (
  ?, ?
) RETURNING *;

-- name: GetBookByTitle :one
SELECT * FROM books
WHERE title = ? LIMIT 1;


-- name: CreateConcreteBook :one
INSERT INTO concrete_books (
  book_id, year_audiobook, parent_path_id
) VALUES (
  ?, ?, ?
) RETURNING *;

-- name: GetConcreteBookByBookAndPath :one
SELECT * FROM concrete_books
WHERE book_id = ? AND parent_path_id = ? LIMIT 1;

