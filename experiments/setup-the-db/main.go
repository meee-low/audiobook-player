package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	sqlc "github.com/meee-low/audiobook-player/internal/db"
	"github.com/meee-low/audiobook-player/sql/migrations"
)

func main() {
	dbPath := "./dev/state/db.db"
	err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	abs, err := filepath.Abs(dbPath)
	log.Printf("Saving to %v", abs)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("ERROR: Could not open the connection to the sqlite database: %v", err)
	}
	defer db.Close()

	log.Printf("Will run migrations")
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Error while migrating: %v", err)
	}

	q := sqlc.New(db)

	res, err := q.CreatePerson(ctx(), "Test Name")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created person with id: %d\n", res.ID)

}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite driver: %w", err)
	}

	d, err := iofs.New(migrations.UpMigrationFiles, "up")
	if err != nil {
		return fmt.Errorf("loading migrations from embed: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("applying migrations: %w", err)
	}

	log.Println("migrations applied successfully")
	return nil
}

func ctx() context.Context {
	return context.Background()
}
