package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/db"
	"github.com/serhappy/rssagg/internal/server"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}
	defer conn.Close()

	app := &app.App{
		DB: db.New(conn),
	}

	srv := &http.Server{
		Handler: server.NewServer(app),
		Addr:    ":" + port,
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
