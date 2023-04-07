package config

import (
	"database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // Import driver for Cloud SQL Proxy
	_ "github.com/lib/pq"
)

func NewInitializedDatabase(config Config) (*sql.DB, error) {
	db, err := NewPostgresSQL(config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresSQL(configuration Config) (*sql.DB, error) {
	username := configuration.Get("DB_USERNAME")
	password := configuration.Get("DB_PASSWORD")
	host := configuration.Get("DB_HOST")
	port := configuration.Get("DB_PORT")
	database := configuration.Get("DB_DATABASE")
	sslMode := configuration.Get("DB_SSL_MODE")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s port=%s sslmode=%s", username, password, database, host, port, sslMode)
	db, err := sql.Open(configuration.Get("DB_CONNECTION"), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
