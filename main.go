package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/menellyn/task-tracker-api/internal/task"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	repo := task.NewPostgresRepository(db)

	tasks, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tasks {
		fmt.Println(t)
	}

}
