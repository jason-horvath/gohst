package main

import (
	"fmt"
	"log"
	"os"
	"time"

	appConfig "gohst/app/config"
	"gohst/internal/config"
	"gohst/internal/db"
	"gohst/internal/migration"
)

func main() {
	// Initialize configuration
	config.InitConfig()
	appConfig.Initialize()      // Initialize app-specific config
	dbConfigs := appConfig.CreateDBConfigs()   // Initialize database configurations
	db.InitDBPool(dbConfigs) // Initialize database connections
	// Initialize database with better error messages for migrations
	if err := db.InitDBForMigrations(); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDBPool()

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		// Initialize the migration model
		migrationModel := migration.NewMigrationModel()
		if err := migrationModel.Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
	case "status":
		// Initialize the migration model
		migrationModel := migration.NewMigrationModel()
		if err := migrationModel.Status(); err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
	case "rollback":
		// Initialize the migration model
		migrationModel := migration.NewMigrationModel()
		if err := migrationModel.Rollback(); err != nil {
			log.Fatal("Rollback failed:", err)
		}
	case "seed":
		// Initialize the seed model
		seedModel := migration.NewSeedModel()
		if err := seedModel.Seed(); err != nil {
			log.Fatal("Seeding failed:", err)
		}
	case "seed:status":
		// Initialize the seed model
		seedModel := migration.NewSeedModel()
		if err := seedModel.SeedStatus(); err != nil {
			log.Fatal("Failed to get seed status:", err)
		}
	case "seed:refresh":
		// Initialize the seed model
		seedModel := migration.NewSeedModel()
		if err := seedModel.SeedRefresh(); err != nil {
			log.Fatal("Seed refresh failed:", err)
		}
	case "seed:rollback":
		// Initialize the seed model
		seedModel := migration.NewSeedModel()
		if err := seedModel.SeedRollback(); err != nil {
			log.Fatal("Seed rollback failed:", err)
		}
	case "full":
		// Run migrations and seeds together
		if err := migration.MigrateAndSeed(); err != nil {
			log.Fatal("Full migration failed:", err)
		}
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate create <migration_name>")
		}
		migrationName := os.Args[2]
		if err := createMigration(migrationName); err != nil {
			log.Fatal("Failed to create migration:", err)
		}
	case "seed:create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate seed:create <seed_name>")
		}
		seedName := os.Args[2]
		if err := createSeed(seedName); err != nil {
			log.Fatal("Failed to create seed:", err)
		}
	default:
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Print(`
Migration Commands:
  run           - Run all pending migrations
  status        - Show migration status
  rollback      - Rollback the last batch of migrations
  seed          - Run all pending seeds
  seed:status   - Show seed status
  seed:refresh  - Clear all seed records and re-run all seeds
  seed:rollback - Rollback the last batch of seeds
  seed:create   - Create a new seed file
  full          - Run migrations and seeds together
  create        - Create a new migration file

Usage:
  migrate run
  migrate status
  migrate rollback
  migrate seed
  migrate seed:status
  migrate seed:refresh
  migrate seed:rollback
  migrate seed:create seed_roles
  migrate full
  migrate create create_users_table
`)
}

func createMigration(name string) error {
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := fmt.Sprintf("database/migrations/%s", filename)

	content := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- Add your migration SQL here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
`, name, time.Now().Format("2006-01-02 15:04:05"))

	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Created migration: %s\n", filepath)
	return nil
}

func createSeed(name string) error {
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := fmt.Sprintf("database/seeds/%s", filename)

	content := fmt.Sprintf(`-- Seed: %s
-- Created: %s

-- Add your seed SQL here
-- Example:
-- INSERT INTO roles (name, description) VALUES
--     ('admin', 'Administrator role'),
--     ('user', 'Regular user role');
`, name, time.Now().Format("2006-01-02 15:04:05"))

	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Created seed: %s\n", filepath)
	return nil
}
