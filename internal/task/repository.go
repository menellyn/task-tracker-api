package task

import (
	"database/sql"
	"errors"
)

var ErrTaskNotFound = errors.New("task not found")

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type PostgresRepository struct {
	db DBTX
}

func NewPostgresRepository(db DBTX) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll() ([]Task, error) {
	tasks := []Task{}

	rows, err := r.db.Query(`
		SELECT id, title, done
		FROM tasks
		ORDER BY id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Done,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *PostgresRepository) GetByID(id int) (Task, error) {
	task := Task{}
	err := r.db.QueryRow(`
		SELECT id, title, done
		FROM tasks
		WHERE id = $1
	`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *PostgresRepository) Add(title string) (Task, error) {
	task := Task{}
	if err := r.db.QueryRow(`
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, done
	`,
		title,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	); err != nil {
		return Task{}, err
	}

	return task, nil

}

func (r *PostgresRepository) MarkDone(id int) error {
	result, err := r.db.Exec(`
		UPDATE tasks
		SET done = true
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(id int) error {
	result, err := r.db.Exec(`
		DELETE FROM tasks
		WHERE id = $1
	`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
