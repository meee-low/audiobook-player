package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2/mp3"
	_ "github.com/mattn/go-sqlite3"

	// "github.com/meee-low/audiobook-player/internal/config"
	"github.com/meee-low/audiobook-player/internal/db"
	"github.com/meee-low/audiobook-player/sql/migrations"
)

func main() {
	SUPPORTED_EXTENSIONS := []string{"mp3"}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <directory>\n", os.Args[0])
	}

	directoryToScan := os.Args[1]
	databasePath := os.Args[2]
	database, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Could not open the database: %v", err)
	}
	defer database.Close()
	err = migrations.RunMigrations(database)
	if err != nil {
		log.Fatalf("Could not apply the migrations: %v", err)
	}

	if _, err := os.Stat(directoryToScan); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", directoryToScan)
	}

	// walk the directory and print all the file names:
	err = filepath.WalkDir(directoryToScan, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Ignore directories, only scan files.
			return nil
		}

		fmt.Printf("File: %s\n", d.Name())
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		fmt.Printf("Abs Path: %s\n", absPath)
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if slices.Contains(SUPPORTED_EXTENSIONS, ext) {
			prepped, err := PrepareFileForDB(path)
			if err != nil {
				return err
			}
			fmt.Printf("Prepped:\n %+v\n", prepped)
			if err := saveToDB(context.Background(), *prepped, database); err != nil {
				return err
			}
			// print_filemetadata(path)
			fmt.Print("\n")
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the directory: %s\n", err)
	}
}

func print_filemetadata(fp string) {
	file, err := os.Open(fp)
	if err != nil {
		log.Fatalf("ERROR: Failed to open file %s: %v", fp, err)
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	printMetadata(metadata)
}

func printMetadata(m tag.Metadata) {
	// Copy-pasted from reference:
	// https://github.com/dhowden/tag/blob/master/cmd/tag/main.go
	fmt.Printf("Metadata Format: %v\n", m.Format())
	fmt.Printf("File Type: %v\n", m.FileType())

	fmt.Printf(" Title: %v\n", m.Title())
	fmt.Printf(" Album: %v\n", m.Album())
	fmt.Printf(" Artist: %v\n", m.Artist())
	fmt.Printf(" Composer: %v\n", m.Composer())
	fmt.Printf(" Genre: %v\n", m.Genre())
	fmt.Printf(" Year: %v\n", m.Year())

	track, trackCount := m.Track()
	fmt.Printf(" Track: %v of %v\n", track, trackCount)

	disc, discCount := m.Disc()
	fmt.Printf(" Disc: %v of %v\n", disc, discCount)

	fmt.Printf(" Picture: %v\n", m.Picture())
	fmt.Printf(" Lyrics: %v\n", m.Lyrics())
	fmt.Printf(" Comment: %v\n", m.Comment())
}

type DatabasePrep struct {
	fileParentPath string
	fileName       string
	fileSizeBytes  int64
	fileExt        string
	fileChecksum   string
	duration       int64
	title          string
	album          string
	artists        []string
	year           int
	track          int
	trackCount     int
	disc           int
	discCount      int
}

func PrepareFileForDB(fp string) (*DatabasePrep, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	checkSum, err := tag.Sum(file)
	if err != nil {
		return nil, err
	}

	duration, err := getMp3Duration(file)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(abs)
	parent := filepath.Dir(abs)
	ext := filepath.Ext(abs)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	track, trackCount := metadata.Track()
	disc, discCount := metadata.Disc()

	result := DatabasePrep{
		fileParentPath: parent,
		fileName:       base,
		fileSizeBytes:  fileSize,
		fileExt:        ext,
		fileChecksum:   checkSum,
		duration:       int64(duration),
		title:          metadata.Title(),
		album:          metadata.Album(),
		artists:        splitArtists(metadata.Artist(), metadata.AlbumArtist()),
		year:           metadata.Year(),
		track:          track,
		trackCount:     trackCount,
		disc:           disc,
		discCount:      discCount,
	}
	return &result, nil
}

func saveToDB(ctx context.Context, dbp DatabasePrep, database *sql.DB) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	q := db.New(tx)

	if err := q.EnsurePath(ctx, dbp.fileParentPath); err != nil {
		tx.Rollback()
		return fmt.Errorf("EnsurePath failed: %w", err)
	}

	path, err := q.GetPathByName(ctx, dbp.fileParentPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("GetPathByName failed: %w", err)
	}

	if err := q.EnsureFileInfo(ctx,
		db.EnsureFileInfoParams{
			Name:         dbp.fileName,
			Extension:    dbp.fileExt,
			Seconds:      dbp.duration,
			ParentPathID: path.ID,
			Checksum:     sql.NullString{String: dbp.fileChecksum, Valid: true},
			Filesize:     sql.NullInt64{Int64: dbp.fileSizeBytes, Valid: dbp.fileSizeBytes > 0},
		}); err != nil {
		tx.Rollback()
		return fmt.Errorf("EnsureFileInfo failed: %w", err)
	}

	file, err := q.GetFileByFullPath(ctx,
		db.GetFileByFullPathParams{
			Name:         dbp.fileName,
			Extension:    dbp.fileExt,
			ParentPathID: path.ID,
		},
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("GetFileByFullPath failed: %w", err)
	}

	if err := q.EnsureBook(ctx,
		db.EnsureBookParams{
			Title:       dbp.album,
			ReleaseYear: sql.NullInt64{Int64: int64(dbp.year), Valid: true}},
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("EnsureBook failed: %w", err)
	}
	book, err := q.GetBookByTitle(ctx, dbp.album)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("GetBookByTitle failed: %w", err)
	}

	if err := q.EnsureConcreteBook(ctx,
		db.EnsureConcreteBookParams{
			BookID:        book.ID,
			YearAudiobook: book.ReleaseYear,
			ParentPathID:  path.ID,
		},
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("EnsureConcreteBook failed: %w", err)
	}

	concreteBook, err := q.GetConcreteBookByBookAndPath(ctx,
		db.GetConcreteBookByBookAndPathParams{
			BookID:       book.ID,
			ParentPathID: path.ID,
		},
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("GetConcreteBookByBookAndPath failed: %w", err)
	}

	if err := q.EnsureAssociateFileToBook(ctx,
		db.EnsureAssociateFileToBookParams{
			ConcreteBookID: concreteBook.ID,
			FileID:         file.ID,
		},
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("EnsureAssociateFileToBook failed: %w", err)
	}

	for _, artist := range dbp.artists {
		if err := q.EnsurePerson(ctx, artist); err != nil {
			tx.Rollback()
			return fmt.Errorf("EnsurePerson failed: %w", err)
		}
		person, err := q.GetPersonByName(ctx, artist)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("GetPersonByName failed: %w", err)
		}
		if err := q.EnsureAssociateAuthorToBook(
			ctx,
			db.EnsureAssociateAuthorToBookParams{AuthorID: person.ID, BookID: book.ID},
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("AssociateAuthorToBook failed: %w", err)
		}
	}

	tx.Commit()
	return nil
}

func splitArtists(artists string, albumArtists string) []string {
	// TODO: use a set
	s := strings.Join([]string{artists, albumArtists}, "%!")
	s = strings.ReplaceAll(artists, ";", "%!")
	s = strings.ReplaceAll(artists, "&", "%!")
	s = strings.ReplaceAll(artists, "/", "%!")

	split := strings.Split(s, "%!")
	result := []string{}
	for _, s := range split {
		result = append(result, strings.TrimSpace(s))
	}

	return result
}

func getMp3Duration(f *os.File) (int, error) {
	return 0, nil
	// This was giving an error, and I really wanted to push this. For now, just returning 0.
	// FIXME: Properly calculate the duration of the mp3.
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return 0, fmt.Errorf("Error while calculating the mp3 duration for %s: %w", f.Name(), err)
	}
	defer streamer.Close()
	duration := streamer.Len() / format.SampleRate.N(time.Second)
	return duration, nil
}
