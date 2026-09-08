package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/menellyn/task-tracker-api/internal/database"
	"github.com/menellyn/task-tracker-api/internal/handler"
	"github.com/menellyn/task-tracker-api/internal/task"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(sqlDB); err != nil {
		log.Fatal(err)
	}

	repo := task.NewORMRepository(db)
	service := task.NewTaskService(repo)
	h := handler.NewTaskHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.GetAll)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
