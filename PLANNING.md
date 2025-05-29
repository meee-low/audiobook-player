This file holds the TODOs/Roadmap

## Basics:
- [x] Play an audio file
- [ ] Pause, skip forward, skip backward etc.
- [ ] Support for multiple audio formats (mp3, m4b, m4a, ...)
- [ ] Scan audiofiles in the chosen directories. Organize them into a database.
- [ ] Store books currently played etc.
- [ ] Display progress (chapters, time etc.)
- [ ] Change playback speed
- [ ] Automatically save position of book and resume from there
- [ ] Edit metadata of books
- [ ] Nice TUI [Charm](https://charm.sh/)
- [ ] Remappable keybinding for actions
- [ ] Multiplatform support (Linux, Windows, MacOS)


## Example of usage/commands
- [ ] `abp`: Starts the TUI audiobook player. If it's not properly configured, runs the setup instead.
- [ ] `abp resume`: Like `abp`, but starts the TUI but automatically resumes the last played book.
- [ ] `abp scan`: Scans the configured directories for audio files and adds them to the database.
- [ ] `abp setup`: Starts the setup wizard to configure the player.
- [ ] `abp config <key> <value>`: Sets a configuration key to that value (will probably just do manual file editing at the start).
- [ ] `abp library`: Opens a webview(?) to edit the metadata more easily.
- [ ] `abp ???`: Some way to interact with a currently running player from the command line. E.g.: play, pause, skip, change speed etc. Maybe use a socket?

## Config options:
- Audio directories to scan
- Audio formats to support
- Autoresume yes/no
- Default playback speed
- Database path (SQLite)
- Keybinds (play, pause, play+pause, background etc.)

## Data format:
- Chapter: title, position, ... (could be multiple files)
- Book <- files, title, author, cover art ... (finish later)
- Series <- collection of books
- Author <- name

Progress: position (chapter and time), last played, maximum read timestamp (so user can go back and resume where they left off)

Need some way to track progress in different chapters. So if it's an anthology, the user can skip around and still track each chapter.


## Pipedreams:
- [ ] Integration with calibre for metadata
- [ ] Cover art if the terminal supports it.

## TUI library view:
### Sorting options:
- All books
    - Order by Author
    - Order by Series
    - Order by BookTitle
    - Order by Genre
    - Order by Read Status
- Select Author -> Order by Series/SeriesNumber/Year
- Select Genre -> Order by Author/Series/SeriesNumber/Year
- Select Series -> Order by SeriesNumber/Year

### Editing:
- Rename authors/books/genre/series/seriesnumber
- mark as read
