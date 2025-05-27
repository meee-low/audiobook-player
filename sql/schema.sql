CREATE TABLE persons (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE books (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER
);

CREATE TABLE authors_books (
    author_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    PRIMARY KEY (author_id, book_id),
    FOREIGN KEY (author_id) REFERENCES persons(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
)

CREATE TABLE narrators_books (
    narrator_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    PRIMARY KEY (narrator_id, book_id),
    FOREIGN KEY (narrator_id) REFERENCES persons(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
)

CREATE TABLE paths (
    id INTEGER PRIMARY KEY,
    path TEXT NOT NULL
)

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    extension TEXT NOT NULL,
    seconds INTEGER NOT NULL,
    parent_path INTEGER NOT NULL,
    FOREIGN KEY (parent_path) REFERENCES paths(id) ON DELETE CASCADE
)

-- Concrete books are the actual instances of books in your hard drive. This lets you have multiple versions of the same book.
CREATE TABLE concrete_books (
    id INTEGER PRIMARY KEY,
    book_id INTEGER NOT NULL,
    year_audiobook INTEGER,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
)

CREATE TABLE books_files (
    id INTEGER PRIMARY KEY,
    order INTEGER,
    concrete_book_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    FOREIGN KEY (concrete_book_id) REFERENCES concrete_books(id) ON DELETE CASCADE
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
)

CREATE TABLE series (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
)

CREATE TABLE series_books (
    series_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    order INTEGER NOT NULL,
    PRIMARY KEY (series_id, book_id),
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
)

-- TODO: Book status/timestamp/progress
