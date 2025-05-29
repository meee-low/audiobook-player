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
);

CREATE TABLE narrators_books (
    narrator_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    PRIMARY KEY (narrator_id, book_id),
    FOREIGN KEY (narrator_id) REFERENCES persons(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE paths (
    id INTEGER PRIMARY KEY,
    path TEXT NOT NULL
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    extension TEXT NOT NULL,
    seconds INTEGER NOT NULL,
    parent_path INTEGER NOT NULL,
    checksum INTEGER NOT NULL,
    filesize INTEGER NOT NULL,
    FOREIGN KEY (parent_path) REFERENCES paths(id) ON DELETE CASCADE
);

-- Concrete books are the actual instances of books in your hard drive. This lets you have multiple versions of the same book.
CREATE TABLE concrete_books (
    id INTEGER PRIMARY KEY,
    book_id INTEGER NOT NULL,
    year_audiobook INTEGER,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE books_files (
    id INTEGER PRIMARY KEY,
    sequence_number INTEGER,
    concrete_book_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    FOREIGN KEY (concrete_book_id) REFERENCES concrete_books(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE series (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE series_books (
    series_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    sequence_number INTEGER,
    PRIMARY KEY (series_id, book_id),
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE file_progress (
    file_id INTEGER NOT NULL,
    seconds INTEGER NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE book_status (
    book_id INTEGER NOT NULL,
    last_played TIMESTAMP,
    status INTEGER NOT NULL,
    PRIMARY KEY (series_id, book_id),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
    FOREIGN KEY (status) REFERENCES status_label(id) ON DELETE CASCADE
);

CREATE TABLE status_label (
    id INTEGER NOT NULL,
    label TEXT NOT NULL UNIQUE
);

INSERT INTO status_label (id, label) VALUES 
    (1, 'not started'),
    (2, 'starting soon'),
    (3, 'did not finish'),
    (4, 'currently reading'),
    (5, 'finished');

