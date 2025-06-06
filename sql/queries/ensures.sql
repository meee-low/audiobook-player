-- name: EnsurePerson :exec
INSERT INTO persons (name)
VALUES (?)
ON CONFLICT(name) DO NOTHING;


-- name: EnsurePath :exec
INSERT INTO paths (path)
VALUES (?)
ON CONFLICT(path) DO NOTHING;


-- name: EnsureFileInfo :exec
INSERT INTO files (
  name, extension, seconds, parent_path_id, checksum, filesize
)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(parent_path_id, name, extension) DO NOTHING;


-- name: EnsureBook :exec
INSERT INTO books (title, release_year)
VALUES (?, ?)
ON CONFLICT(title) DO NOTHING;


-- name: EnsureConcreteBook :exec
INSERT INTO concrete_books (book_id, year_audiobook, parent_path_id)
VALUES (?, ?, ?)
ON CONFLICT(book_id, parent_path_id) DO NOTHING;



-- name: EnsureAssociateFileToBook :exec
INSERT INTO books_files (concrete_book_id, file_id, sequence_number)
VALUES ( ?, ?, ?)
ON CONFLICT(concrete_book_id, file_id) DO NOTHING;


-- name: EnsureAssociateAuthorToBook :exec
INSERT INTO authors_books (author_id, book_id)
VALUES ( ?, ? )
ON CONFLICT(author_id, book_id) DO NOTHING;
