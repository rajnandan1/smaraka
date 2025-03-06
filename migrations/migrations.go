package migrations

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rajnandan1/smaraka/postgres"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

func DoPostgresMigrationsUp(connString string) {
	m, err := migrate.New(
		"file://migrations",
		connString)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	// Run migrations up to the latest version
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}

func DoPostgresMigrationsDown(connString string) {
	m, err := migrate.New(
		"file://migrations",
		connString)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	// Run migrations down to the latest version
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}

func DoRiverMigrationUp(dbPool postgres.Postgres) {
	driver := riverpgxv5.New(dbPool.GetConnectionPool())
	migrator, err := rivermigrate.New(driver, nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, riverErr := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if riverErr != nil {
		panic(riverErr)
	}
	fmt.Println("Migrations for River applied successfully!")
}

func DoRiverMigrationDown(dbPool postgres.Postgres) {
	driver := riverpgxv5.New(dbPool.GetConnectionPool())
	migrator, err := rivermigrate.New(driver, nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, riverErr := migrator.Migrate(ctx, rivermigrate.DirectionDown, nil)
	if riverErr != nil {
		panic(riverErr)
	}
	fmt.Println("Migrations for River applied successfully!")
}
