package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gohst/internal/models"
)

type Migration struct {
    ID        int       `db:"id"`
    Migration string    `db:"migration"`
    Batch     int       `db:"batch"`
    RunAt     time.Time `db:"run_at"`
}

type Seed struct {
    ID     int       `db:"id"`
    Seed   string    `db:"seed"`
    Batch  int       `db:"batch"`
    RunAt  time.Time `db:"run_at"`
}

type MigrationFile struct {
    Filename string
    Path     string
    Content  string
}

type SeedFile struct {
    Filename string
    Path     string
    Content  string
}

type MigrationModel struct {
    *models.Model[Migration]
}

type SeedModel struct {
    *models.Model[Seed]
}

// NewMigrationModel creates a new migration model instance
func NewMigrationModel() *MigrationModel {
    return &MigrationModel{
        Model: models.NewModel[Migration]("migrations"),
    }
}

// NewSeedModel creates a new seed model instance
func NewSeedModel() *SeedModel {
    return &SeedModel{
        Model: models.NewModel[Seed]("seeds"),
    }
}

// CreateMigrationsTable creates the migrations tracking table
func (m *MigrationModel) CreateMigrationsTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        migration VARCHAR(255) NOT NULL,
        batch INTEGER NOT NULL,
        run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

    _, err := m.GetDB().Exec(query)
    if err != nil {
        return fmt.Errorf("failed to create migrations table: %v", err)
    }

    return nil
}

// GetRunMigrations returns all migrations that have been run
func (m *MigrationModel) GetRunMigrations() ([]*Migration, error) {
    var migrations []*Migration

    rows, err := m.GetDB().Query("SELECT id, migration, batch, run_at FROM migrations ORDER BY batch ASC, id ASC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        migration := &Migration{}
        err := rows.Scan(&migration.ID, &migration.Migration, &migration.Batch, &migration.RunAt)
        if err != nil {
            return nil, err
        }
        migrations = append(migrations, migration)
    }

    return migrations, nil
}

// GetMigrationFiles returns all migration files from the migrations directory
func (m *MigrationModel) GetMigrationFiles() ([]MigrationFile, error) {
    migrationDir := "database/migrations"
    var files []MigrationFile

    err := filepath.Walk(migrationDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
            content, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            files = append(files, MigrationFile{
                Filename: info.Name(),
                Path:     path,
                Content:  string(content),
            })
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Sort files by filename (which should include timestamp)
    sort.Slice(files, func(i, j int) bool {
        return files[i].Filename < files[j].Filename
    })

    return files, nil
}

// GetPendingMigrations returns migrations that haven't been run yet
func (m *MigrationModel) GetPendingMigrations() ([]MigrationFile, error) {
    allFiles, err := m.GetMigrationFiles()
    if err != nil {
        return nil, err
    }

    runMigrations, err := m.GetRunMigrations()
    if err != nil {
        return nil, err
    }

    // Create a map of run migrations for quick lookup
    runMap := make(map[string]bool)
    for _, migration := range runMigrations {
        runMap[migration.Migration] = true
    }

    var pending []MigrationFile
    for _, file := range allFiles {
        if !runMap[file.Filename] {
            pending = append(pending, file)
        }
    }

    return pending, nil
}

// RunMigration executes a single migration
func (m *MigrationModel) RunMigration(migrationFile MigrationFile, batch int) error {
    // Execute the migration SQL
    _, err := m.GetDB().Exec(migrationFile.Content)
    if err != nil {
        return fmt.Errorf("failed to execute migration %s: %v", migrationFile.Filename, err)
    }

    // Record the migration as run
    _, err = m.GetDB().Exec("INSERT INTO migrations (migration, batch) VALUES ($1, $2)", migrationFile.Filename, batch)
    if err != nil {
        return fmt.Errorf("failed to record migration %s: %v", migrationFile.Filename, err)
    }

    return nil
}

// GetNextBatch returns the next batch number
func (m *MigrationModel) GetNextBatch() (int, error) {
    var batch int
    err := m.GetDB().QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM migrations").Scan(&batch)
    if err != nil {
        return 0, err
    }
    return batch, nil
}

// Migrate runs all pending migrations
func (m *MigrationModel) Migrate() error {
    if err := m.CreateMigrationsTable(); err != nil {
        return err
    }

    pending, err := m.GetPendingMigrations()
    if err != nil {
        return err
    }

    if len(pending) == 0 {
        log.Println("No pending migrations to run")
        return nil
    }

    batch, err := m.GetNextBatch()
    if err != nil {
        return err
    }

    // Start a transaction
    tx, err := m.GetDB().Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    log.Printf("Running %d migrations in batch %d", len(pending), batch)

    for _, migrationFile := range pending {
        log.Printf("Running migration: %s", migrationFile.Filename)

        if err := m.RunMigration(migrationFile, batch); err != nil {
            return err
        }
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        return err
    }

    log.Printf("Successfully ran %d migrations", len(pending))
    return nil
}

// Refresh drops all tables and re-runs all migrations
func (m *MigrationModel) Refresh() error {
    // Get all table names
    rows, err := m.GetDB().Query(`
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public'
        AND table_type = 'BASE TABLE'
    `)
    if err != nil {
        return fmt.Errorf("failed to get tables: %v", err)
    }
    defer rows.Close()

    var tables []string
    for rows.Next() {
        var table string
        if err := rows.Scan(&table); err != nil {
            return fmt.Errorf("failed to scan table name: %v", err)
        }
        tables = append(tables, table)
    }

    if len(tables) == 0 {
        log.Println("No tables to drop")
    } else {
        log.Printf("Dropping %d tables...", len(tables))

        // Disable foreign key checks temporarily if needed, or just use CASCADE
        for _, table := range tables {
            log.Printf("Dropping table: %s", table)
            _, err := m.GetDB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS \"%s\" CASCADE", table))
            if err != nil {
                return fmt.Errorf("failed to drop table %s: %v", table, err)
            }
        }
        log.Println("Successfully dropped all tables")
    }

    // Re-run migrations
    log.Println("Re-running all migrations...")
    return m.Migrate()
}

// Status shows migration status
func (m *MigrationModel) Status() error {
    if err := m.CreateMigrationsTable(); err != nil {
        return err
    }

    allFiles, err := m.GetMigrationFiles()
    if err != nil {
        return err
    }

    runMigrations, err := m.GetRunMigrations()
    if err != nil {
        return err
    }

    // Create a map of run migrations for quick lookup
    runMap := make(map[string]*Migration)
    for _, migration := range runMigrations {
        runMap[migration.Migration] = migration
    }

    fmt.Println("\n=== Migration Status ===")
    fmt.Printf("%-50s %-10s %-10s %-20s\n", "Migration", "Status", "Batch", "Run At")
    fmt.Println(strings.Repeat("-", 90))

    for _, file := range allFiles {
        if migration, exists := runMap[file.Filename]; exists {
            fmt.Printf("%-50s %-10s %-10d %-20s\n",
                file.Filename,
                "✅ RUN",
                migration.Batch,
                migration.RunAt.Format("2006-01-02 15:04:05"))
        } else {
            fmt.Printf("%-50s %-10s %-10s %-20s\n",
                file.Filename,
                "❌ PENDING",
                "-",
                "-")
        }
    }

    pending, _ := m.GetPendingMigrations()
    fmt.Printf("\nPending migrations: %d\n", len(pending))

    return nil
}

// ============ SEED FUNCTIONALITY ============

// CreateSeedsTable creates the seeds tracking table
func (s *SeedModel) CreateSeedsTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS seeds (
        id SERIAL PRIMARY KEY,
        seed VARCHAR(255) NOT NULL,
        batch INTEGER NOT NULL,
        run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

    _, err := s.GetDB().Exec(query)
    if err != nil {
        return fmt.Errorf("failed to create seeds table: %v", err)
    }

    return nil
}

// GetRunSeeds returns all seeds that have been run
func (s *SeedModel) GetRunSeeds() ([]*Seed, error) {
    var seeds []*Seed

    rows, err := s.GetDB().Query("SELECT id, seed, batch, run_at FROM seeds ORDER BY batch ASC, id ASC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        seed := &Seed{}
        err := rows.Scan(&seed.ID, &seed.Seed, &seed.Batch, &seed.RunAt)
        if err != nil {
            return nil, err
        }
        seeds = append(seeds, seed)
    }

    return seeds, nil
}

// GetSeedFiles returns all seed files from the seeds directory
func (s *SeedModel) GetSeedFiles() ([]SeedFile, error) {
    seedDir := "database/seeds"
    var files []SeedFile

    err := filepath.Walk(seedDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
            content, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            files = append(files, SeedFile{
                Filename: info.Name(),
                Path:     path,
                Content:  string(content),
            })
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Sort files by filename (which should include timestamp)
    sort.Slice(files, func(i, j int) bool {
        return files[i].Filename < files[j].Filename
    })

    return files, nil
}

// GetPendingSeeds returns seeds that haven't been run yet
func (s *SeedModel) GetPendingSeeds() ([]SeedFile, error) {
    allFiles, err := s.GetSeedFiles()
    if err != nil {
        return nil, err
    }

    runSeeds, err := s.GetRunSeeds()
    if err != nil {
        return nil, err
    }

    // Create a map of run seeds for quick lookup
    runMap := make(map[string]bool)
    for _, seed := range runSeeds {
        runMap[seed.Seed] = true
    }

    var pending []SeedFile
    for _, file := range allFiles {
        if !runMap[file.Filename] {
            pending = append(pending, file)
        }
    }

    return pending, nil
}

// GetNextSeedBatch returns the next batch number for seeds
func (s *SeedModel) GetNextSeedBatch() (int, error) {
    var batch int
    err := s.GetDB().QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM seeds").Scan(&batch)
    if err != nil {
        return 0, err
    }
    return batch, nil
}

// RunSeed executes a single seed
func (s *SeedModel) RunSeed(seedFile SeedFile, batch int) error {
    // Execute the seed SQL
    _, err := s.GetDB().Exec(seedFile.Content)
    if err != nil {
        return fmt.Errorf("failed to execute seed %s: %v", seedFile.Filename, err)
    }

    // Record the seed as run
    _, err = s.GetDB().Exec("INSERT INTO seeds (seed, batch) VALUES ($1, $2)", seedFile.Filename, batch)
    if err != nil {
        return fmt.Errorf("failed to record seed %s: %v", seedFile.Filename, err)
    }

    return nil
}

// Seed runs all pending seeds
func (s *SeedModel) Seed() error {
    if err := s.CreateSeedsTable(); err != nil {
        return err
    }

    pending, err := s.GetPendingSeeds()
    if err != nil {
        return err
    }

    if len(pending) == 0 {
        log.Println("No pending seeds to run")
        return nil
    }

    batch, err := s.GetNextSeedBatch()
    if err != nil {
        return err
    }

    // Start a transaction
    tx, err := s.GetDB().Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    log.Printf("Running %d seeds", len(pending))

    for _, seedFile := range pending {
        log.Printf("Running seed: %s", seedFile.Filename)

        if err := s.RunSeed(seedFile, batch); err != nil {
            return err
        }
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        return err
    }

    log.Printf("Successfully ran %d seeds", len(pending))
    return nil
}

// SeedStatus shows seed status
func (s *SeedModel) SeedStatus() error {
    if err := s.CreateSeedsTable(); err != nil {
        return err
    }

    allFiles, err := s.GetSeedFiles()
    if err != nil {
        return err
    }

    runSeeds, err := s.GetRunSeeds()
    if err != nil {
        return err
    }

    // Create a map of run seeds for quick lookup
    runMap := make(map[string]*Seed)
    for _, seed := range runSeeds {
        runMap[seed.Seed] = seed
    }

    fmt.Println("\n=== Seed Status ===")
    fmt.Printf("%-50s %-10s %-10s %-20s\n", "Seed", "Status", "Batch", "Run At")
    fmt.Println(strings.Repeat("-", 90))

    for _, file := range allFiles {
        if seed, exists := runMap[file.Filename]; exists {
            fmt.Printf("%-50s %-10s %-10d %-20s\n",
                file.Filename,
                "✅ RUN",
                seed.Batch,
                seed.RunAt.Format("2006-01-02 15:04:05"))
        } else {
            fmt.Printf("%-50s %-10s %-10s %-20s\n",
                file.Filename,
                "❌ PENDING",
                "-",
                "-")
        }
    }

    pending, _ := s.GetPendingSeeds()
    fmt.Printf("\nPending seeds: %d\n", len(pending))

    return nil
}

// SeedRefresh drops all seed records and re-runs all seeds
func (s *SeedModel) SeedRefresh() error {
    if err := s.CreateSeedsTable(); err != nil {
        return err
    }

    log.Println("Clearing seed records...")
    _, err := s.GetDB().Exec("DELETE FROM seeds")
    if err != nil {
        return fmt.Errorf("failed to clear seed records: %v", err)
    }

    log.Println("Re-running all seeds...")
    return s.Seed()
}

// SeedRollback rolls back the last batch of seeds
func (s *SeedModel) SeedRollback() error {
    // Get the last batch number
    var lastBatch int
    err := s.GetDB().QueryRow("SELECT COALESCE(MAX(batch), 0) FROM seeds").Scan(&lastBatch)
    if err != nil {
        return err
    }

    if lastBatch == 0 {
        log.Println("No seeds to rollback")
        return nil
    }

    // Get seeds from the last batch
    rows, err := s.GetDB().Query("SELECT seed FROM seeds WHERE batch = $1 ORDER BY id DESC", lastBatch)
    if err != nil {
        return err
    }
    defer rows.Close()

    var seeds []string
    for rows.Next() {
        var seed string
        if err := rows.Scan(&seed); err != nil {
            return err
        }
        seeds = append(seeds, seed)
    }

    log.Printf("Rolling back %d seeds from batch %d", len(seeds), lastBatch)
    log.Println("⚠️  Note: This only removes seed records, not the actual data inserted by seeds")

    // Delete the seed records
    _, err = s.GetDB().Exec("DELETE FROM seeds WHERE batch = $1", lastBatch)
    if err != nil {
        return err
    }

    log.Printf("Successfully rolled back seed batch %d", lastBatch)
    return nil
}

// ============ COMBINED FUNCTIONALITY ============

// MigrateAndSeed runs migrations first, then seeds
func MigrateAndSeed() error {
    log.Println("🚀 Starting full database setup...")

    // Run migrations first
    log.Println("📋 Running migrations...")
    migrationModel := NewMigrationModel()
    if err := migrationModel.Migrate(); err != nil {
        return fmt.Errorf("migration failed: %v", err)
    }

    // Run seeds after migrations
    log.Println("🌱 Running seeds...")
    seedModel := NewSeedModel()
    if err := seedModel.Seed(); err != nil {
        return fmt.Errorf("seeding failed: %v", err)
    }

    log.Println("✅ Database setup complete!")
    return nil
}

// Rollback rolls back the last batch of migrations
func (m *MigrationModel) Rollback() error {
    // Get the last batch number
    var lastBatch int
    err := m.GetDB().QueryRow("SELECT COALESCE(MAX(batch), 0) FROM migrations").Scan(&lastBatch)
    if err != nil {
        return err
    }

    if lastBatch == 0 {
        log.Println("No migrations to rollback")
        return nil
    }

    // Get migrations from the last batch
    rows, err := m.GetDB().Query("SELECT migration FROM migrations WHERE batch = $1 ORDER BY id DESC", lastBatch)
    if err != nil {
        return err
    }
    defer rows.Close()

    var migrations []string
    for rows.Next() {
        var migration string
        if err := rows.Scan(&migration); err != nil {
            return err
        }
        migrations = append(migrations, migration)
    }

    log.Printf("Rolling back %d migrations from batch %d", len(migrations), lastBatch)

    // Delete the migration records
    _, err = m.GetDB().Exec("DELETE FROM migrations WHERE batch = $1", lastBatch)
    if err != nil {
        return err
    }

    log.Printf("Successfully rolled back batch %d", lastBatch)
    return nil
}

// RollbackSeeds rolls back the last batch of seeds
func (s *SeedModel) RollbackSeeds() error {
    // Get the last seed ID
    var lastSeedID int
    err := s.GetDB().QueryRow("SELECT COALESCE(MAX(id), 0) FROM seeds").Scan(&lastSeedID)
    if err != nil {
        return err
    }

    if lastSeedID == 0 {
        log.Println("No seeds to rollback")
        return nil
    }

    // Get seeds from the last batch
    rows, err := s.GetDB().Query("SELECT seed FROM seeds ORDER BY id DESC")
    if err != nil {
        return err
    }
    defer rows.Close()

    var seeds []string
    for rows.Next() {
        var seed string
        if err := rows.Scan(&seed); err != nil {
            return err
        }
        seeds = append(seeds, seed)
    }

    log.Printf("Rolling back %d seeds", len(seeds))

    // Delete the seed records
    _, err = s.GetDB().Exec("DELETE FROM seeds WHERE id <= $1", lastSeedID)
    if err != nil {
        return err
    }

    log.Printf("Successfully rolled back seeds up to ID %d", lastSeedID)
    return nil
}
