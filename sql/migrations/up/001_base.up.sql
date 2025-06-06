CREATE TABLE persons (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    UNIQUE (name)
);

CREATE TABLE aliases (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE persons_aliases (
    person_id INTEGER NOT NULL,
    alias_id INTEGER NOT NULL,
    PRIMARY KEY (person_id, alias_id),
    FOREIGN KEY (person_id) REFERENCES persons(id) ON DELETE CASCADE,
    FOREIGN KEY (alias_id) REFERENCES aliases(id) ON DELETE CASCADE
);

CREATE TABLE books (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER,
    UNIQUE (title)
);

CREATE TABLE books_aliases (
    book_id INTEGER NOT NULL,
    alias_id INTEGER NOT NULL,
    PRIMARY KEY (book_id, alias_id),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (alias_id) REFERENCES aliases(id) ON DELETE CASCADE
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
    path TEXT NOT NULL,
    UNIQUE (path)
);

CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    extension TEXT NOT NULL,
    seconds INTEGER NOT NULL,
    parent_path_id INTEGER NOT NULL,
    checksum TEXT,
    filesize INTEGER,
    FOREIGN KEY (parent_path_id) REFERENCES paths(id) ON DELETE CASCADE,
    UNIQUE (name, extension, parent_path_id)
);

-- Concrete books are the actual instances of books in your hard drive. This lets you have multiple versions of the same book.
CREATE TABLE concrete_books (
    id INTEGER PRIMARY KEY,
    book_id INTEGER NOT NULL,
    parent_path_id INTEGER NOT NULL,
    year_audiobook INTEGER,
    UNIQUE (book_id, parent_path_id),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_path_id) REFERENCES paths(id) ON DELETE CASCADE
);

CREATE TABLE books_files (
    id INTEGER PRIMARY KEY,
    sequence_number INTEGER,
    concrete_book_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    UNIQUE (concrete_book_id, file_id),
    FOREIGN KEY (concrete_book_id) REFERENCES concrete_books(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE series (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    UNIQUE (name)
);

CREATE TABLE series_aliases (
    series_id INTEGER NOT NULL,
    alias_id INTEGER NOT NULL,
    PRIMARY KEY (series_id, alias_id),
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    FOREIGN KEY (alias_id) REFERENCES aliases(id) ON DELETE CASCADE
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
    PRIMARY KEY (file_id),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE book_status (
    book_id INTEGER NOT NULL,
    last_played TIMESTAMP,
    status INTEGER NOT NULL,
    PRIMARY KEY (book_id),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
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

