package config

import (
	//_ "cloud.google.com/go/cloudsqlconn"
	"database/sql"
	"fmt"
	//_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	//_ "github.com/jackc/pgx/v4/stdlib"
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, database, sslMode)
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
