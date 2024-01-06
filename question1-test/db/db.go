package db

import (
	"context"
	"io/ioutil"
	"os"

	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	dbUrl := os.Getenv("DATABASE_URL") // Set this environment variable
	dbPool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	// Initialize tables from init.sql
	if err := initializeTables(dbPool); err != nil {
		return nil, err
	}

	return dbPool, nil
}

func initializeTables(dbPool *pgxpool.Pool) error {
	// Adjust the path to the location of your init.sql file
	initSQL, err := ioutil.ReadFile("db/init.sql")
	if err != nil {
		log.Printf("Error reading init.sql file: %v", err)
		return err
	}

	_, err = dbPool.Exec(context.Background(), string(initSQL))
	if err != nil {
		log.Printf("Error executing init.sql commands: %v", err)
		return err
	}

	log.Println("Database tables initialized successfully")
	return nil
}
