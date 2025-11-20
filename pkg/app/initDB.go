package app

import (
	"context"
	"database/sql"
	"log"
	"os"

	DB "github.com/egor_lukyanovich/avito/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	Queries *DB.Queries
	DB      *sql.DB
}

func InitDB() (Storage, string) {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error to load .env in app.go:", err)
	}

	dataBaseUrl := os.Getenv("DATABASE_URL")
	if dataBaseUrl == "" {
		log.Fatal("DataBaseUrl is not found in .env:")
	}

	// port := os.Getenv("SERVER_PORT")
	// if port == "" {
	// 	log.Fatal("SERVER_PORT is not found in .env:")
	// }

	db, err := sql.Open("pgx", dataBaseUrl)
	if err != nil {
		log.Fatal("Failed to open a database:", err)
	}

	err = goose.Up(db, "/app/sql/schema")
	if err != nil {
		log.Println("Failed to up a database with goose:", err)
	}

	db.PingContext(context.Background())

	queries := DB.New(db)

	storage := Storage{
		Queries: queries,
		DB:      db,
	}

	port := "" // заглушка
	return storage, port
}
