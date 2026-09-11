package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	h := handler.NewTaskHandler(service, validate)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.GetAll)
	mux.HandleFunc("GET /tasks/actual", h.GetActual)
	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	mux.HandleFunc("POST /tasks", h.Create)
	mux.HandleFunc("PUT /tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)
	mux.HandleFunc("PATCH /tasks/{id}/done", h.MarkDone)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
