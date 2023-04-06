package config

import (
	"database/sql"
	"fmt"
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

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", username, password, host, port, database, sslMode)
	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
