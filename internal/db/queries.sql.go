// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const createBook = `-- name: CreateBook :one
INSERT INTO books (
  title, release_year
) VALUES (
  ?, ?
) RETURNING id, title, release_year
`

type CreateBookParams struct {
	Title       string
	ReleaseYear sql.NullInt64
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, createBook, arg.Title, arg.ReleaseYear)
	var i Book
	err := row.Scan(&i.ID, &i.Title, &i.ReleaseYear)
	return i, err
}

const createConcreteBook = `-- name: CreateConcreteBook :one
INSERT INTO concrete_books (
  book_id, year_audiobook, parent_path_id
) VALUES (
  ?, ?, ?
) RETURNING id, book_id, parent_path_id, year_audiobook
`

type CreateConcreteBookParams struct {
	BookID        int64
	YearAudiobook sql.NullInt64
	ParentPathID  int64
}

func (q *Queries) CreateConcreteBook(ctx context.Context, arg CreateConcreteBookParams) (ConcreteBook, error) {
	row := q.db.QueryRowContext(ctx, createConcreteBook, arg.BookID, arg.YearAudiobook, arg.ParentPathID)
	var i ConcreteBook
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.ParentPathID,
		&i.YearAudiobook,
	)
	return i, err
}

const createFileInfo = `-- name: CreateFileInfo :one
INSERT INTO files (
  name, extension, seconds, parent_path_id, checksum, filesize
) VALUES (
  ?, ?, ?, ?, ?, ?
) RETURNING id, name, extension, seconds, parent_path_id, checksum, filesize
`

type CreateFileInfoParams struct {
	Name         string
	Extension    string
	Seconds      int64
	ParentPathID int64
	Checksum     sql.NullString
	Filesize     sql.NullInt64
}

func (q *Queries) CreateFileInfo(ctx context.Context, arg CreateFileInfoParams) (File, error) {
	row := q.db.QueryRowContext(ctx, createFileInfo,
		arg.Name,
		arg.Extension,
		arg.Seconds,
		arg.ParentPathID,
		arg.Checksum,
		arg.Filesize,
	)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Extension,
		&i.Seconds,
		&i.ParentPathID,
		&i.Checksum,
		&i.Filesize,
	)
	return i, err
}

const createParentPath = `-- name: CreateParentPath :one
INSERT INTO paths ( path )
VALUES ( ? )
RETURNING id, path
`

func (q *Queries) CreateParentPath(ctx context.Context, path string) (Path, error) {
	row := q.db.QueryRowContext(ctx, createParentPath, path)
	var i Path
	err := row.Scan(&i.ID, &i.Path)
	return i, err
}

const createPerson = `-- name: CreatePerson :one
INSERT INTO persons (
  name
) VALUES (
  ?
) RETURNING id, name
`

func (q *Queries) CreatePerson(ctx context.Context, name string) (Person, error) {
	row := q.db.QueryRowContext(ctx, createPerson, name)
	var i Person
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getBookByTitle = `-- name: GetBookByTitle :one
SELECT id, title, release_year FROM books
WHERE title = ? LIMIT 1
`

func (q *Queries) GetBookByTitle(ctx context.Context, title string) (Book, error) {
	row := q.db.QueryRowContext(ctx, getBookByTitle, title)
	var i Book
	err := row.Scan(&i.ID, &i.Title, &i.ReleaseYear)
	return i, err
}

const getConcreteBookByBookAndPath = `-- name: GetConcreteBookByBookAndPath :one
SELECT id, book_id, parent_path_id, year_audiobook FROM concrete_books
WHERE book_id = ? AND parent_path_id = ? LIMIT 1
`

type GetConcreteBookByBookAndPathParams struct {
	BookID       int64
	ParentPathID int64
}

func (q *Queries) GetConcreteBookByBookAndPath(ctx context.Context, arg GetConcreteBookByBookAndPathParams) (ConcreteBook, error) {
	row := q.db.QueryRowContext(ctx, getConcreteBookByBookAndPath, arg.BookID, arg.ParentPathID)
	var i ConcreteBook
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.ParentPathID,
		&i.YearAudiobook,
	)
	return i, err
}

const getFileByFullPath = `-- name: GetFileByFullPath :one
SELECT id, name, extension, seconds, parent_path_id, checksum, filesize FROM files
WHERE (
  parent_path_id = ? AND
  name = ? AND
  extension = ?
) LIMIT 1
`

type GetFileByFullPathParams struct {
	ParentPathID int64
	Name         string
	Extension    string
}

func (q *Queries) GetFileByFullPath(ctx context.Context, arg GetFileByFullPathParams) (File, error) {
	row := q.db.QueryRowContext(ctx, getFileByFullPath, arg.ParentPathID, arg.Name, arg.Extension)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Extension,
		&i.Seconds,
		&i.ParentPathID,
		&i.Checksum,
		&i.Filesize,
	)
	return i, err
}

const getPathByName = `-- name: GetPathByName :one
SELECT id, path FROM paths
WHERE path = ? LIMIT 1
`

func (q *Queries) GetPathByName(ctx context.Context, path string) (Path, error) {
	row := q.db.QueryRowContext(ctx, getPathByName, path)
	var i Path
	err := row.Scan(&i.ID, &i.Path)
	return i, err
}

const getPerson = `-- name: GetPerson :one
SELECT id, name FROM persons
WHERE id = ? LIMIT 1
`

func (q *Queries) GetPerson(ctx context.Context, id int64) (Person, error) {
	row := q.db.QueryRowContext(ctx, getPerson, id)
	var i Person
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getPersonByName = `-- name: GetPersonByName :one
SELECT id, name FROM persons
WHERE name = ? LIMIT 1
`

func (q *Queries) GetPersonByName(ctx context.Context, name string) (Person, error) {
	row := q.db.QueryRowContext(ctx, getPersonByName, name)
	var i Person
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
