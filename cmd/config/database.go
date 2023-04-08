package config

import (
	_ "cloud.google.com/go/cloudsqlconn"
	"database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewInitializedDatabase(config Config) (*sql.DB, error) {
	db, err := NewPostgresSQL(config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresSQL(conn Config) (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_unix.go: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser         = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PASS") // e.g. 'my-db-password'
		unixSocketPath = mustGetenv("DB_HOST") // e.g. '/cloudsql/project:region:instance'
		dbName         = mustGetenv("DB_NAME") // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// ...

	return dbPool, nil
}
func MigrateDatabase() error {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	instanceSocket := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Build the connection string
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		dbUser, dbPass, dbName, instanceSocket)

	// Open a connection to the database
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Set the migration path
	migrationPath := "file://database/postgres/migration"
	// Run migrations
	sourceDriver, err := (&file.File{}).Open(migrationPath)
	if err != nil {
		return err
	}
	defer sourceDriver.Close()

	m, err := migrate.NewWithSourceInstance("file", sourceDriver, connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations completed successfully")

	return nil
}
