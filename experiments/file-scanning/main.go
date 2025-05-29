package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	// "github.com/meee-low/audiobook-player/internal/config"
	"github.com/dhowden/tag"
)

func main() {
	SUPPORTED_EXTENSIONS := []string{"mp3"}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <directory>\n", os.Args[0])
	}

	directory := os.Args[1]

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", directory)
	}

	// walk the directory and print all the file names:
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
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
			print_filemetadata(path)
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
	title          string
	album          string
	artists        []string
	year           int
	track          int
	trackCount     int
	disc           int
	discCount      int
}

func BlankDBPrep() DatabasePrep {
	return DatabasePrep{}
}

func FileToDB(fp string) (DatabasePrep, error) {
	file, err := os.Open(fp)
	if err != nil {
		return BlankDBPrep(), err
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return BlankDBPrep(), err
	}
	checkSum, err := tag.Sum(file)
	if err != nil {
		return BlankDBPrep(), err
	}

	abs, err := filepath.Abs(fp)
	if err != nil {
		return BlankDBPrep(), err
	}

	base := filepath.Base(abs)
	parent := filepath.Dir(abs)
	ext := filepath.Ext(abs)
	fileInfo, err := file.Stat()
	if err != nil {
		return BlankDBPrep(), err
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
		title:          metadata.Title(),
		album:          metadata.Album(),
		artists:        splitArtists(metadata.Artist(), metadata.AlbumArtist()),
		year:           metadata.Year(),
		track:          track,
		trackCount:     trackCount,
		disc:           disc,
		discCount:      discCount,
	}
	return result, nil
}

func saveToDB(dbp DatabasePrep) {
	/**
		dbfileinfo <- fileParentPath, fileName, fileSizeBytes, fileExt, fileChecksum
	    book <- title, year
	    concrete_book <- ...
	    persons <- artists
	    joins...
	**/
}

func splitArtists(artists string, albumArtists string) []string {
	// TODO: use a set and use albumArtists
	s := strings.Join([]string{artists, albumArtists}, "%!")
	s = strings.ReplaceAll(artists, ";", "%!")
	s = strings.ReplaceAll(artists, "&", "%!")

	return strings.Split(s, "%!")
}
