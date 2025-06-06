package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/meee-low/audiobook-player/internal/db"
	"github.com/meee-low/audiobook-player/sql/migrations"
)

func main() {
	dbPath := "./dev/state/db.db"
	err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	abs, err := filepath.Abs(dbPath)
	log.Printf("Saving to %v", abs)

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("ERROR: Could not open the connection to the sqlite database: %v", err)
	}
	defer database.Close()

	log.Printf("Will run migrations")
	err = migrations.RunMigrations(database)
	if err != nil {
		log.Fatalf("Error while migrating: %v", err)
	}

	q := db.New(database)

	res, err := q.CreatePerson(ctx(), "Test Name")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created person with id: %d\n", res.ID)

}

func ctx() context.Context {
	return context.Background()
}
