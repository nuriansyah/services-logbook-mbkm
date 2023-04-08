package config

import (
	_ "cloud.google.com/go/cloudsqlconn"
	"database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
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
