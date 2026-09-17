package task

import (
	"database/sql"
	"errors"
	"time"
)

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
		SELECT id, title, description, schedule_date, deadline, done
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
			&task.Description,
			&task.ScheduleDate,
			&task.Deadline,
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

func (r *PostgresRepository) GetActual() ([]Task, error) {
	tasks := []Task{}

	rows, err := r.db.Query(`
		SELECT id, title, description, schedule_date, deadline, done
		FROM tasks
		WHERE done = false
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
			&task.Description,
			&task.ScheduleDate,
			&task.Deadline,
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
		SELECT id, title, description, schedule_date, deadline, done
		FROM tasks
		WHERE id = $1
	`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.ScheduleDate,
		&task.Deadline,
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

func (r *PostgresRepository) Add(task Task) (Task, error) {
	if err := r.db.QueryRow(`
		INSERT INTO tasks (title, description, schedule_date, deadline)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, schedule_date, deadline, done
	`,
		task.Title,
		task.Description,
		task.ScheduleDate,
		task.Deadline,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.ScheduleDate,
		&task.Deadline,
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

func (r *PostgresRepository) Update(id int, updateData map[string]interface{}) (Task, error) {
	currentTask := Task{}

	err := r.db.QueryRow(`
        SELECT id, title, description, schedule_date, deadline, done
        FROM tasks
        WHERE id = $1
    `, id).Scan(
		&currentTask.ID,
		&currentTask.Title,
		&currentTask.Description,
		&currentTask.ScheduleDate,
		&currentTask.Deadline,
		&currentTask.Done,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}

	for key, value := range updateData {
		switch key {
		case "title":
			currentTask.Title = value.(string)

		case "description":
			currentTask.Description = value.(*string)

		case "schedule_date":
			currentTask.ScheduleDate = value.(*time.Time)

		case "deadline":
			currentTask.Deadline = value.(*time.Time)

		case "done":
			currentTask.Done = value.(bool)
		}
	}

	err = r.db.QueryRow(`
        UPDATE tasks
        SET title = $1,
            description = $2,
            schedule_date = $3,
            deadline = $4,
            done = $5
        WHERE id = $6
        RETURNING id, title, description, schedule_date, deadline, done
    `,
		currentTask.Title,
		currentTask.Description,
		currentTask.ScheduleDate,
		currentTask.Deadline,
		currentTask.Done,
		currentTask.ID,
	).Scan(
		&currentTask.ID,
		&currentTask.Title,
		&currentTask.Description,
		&currentTask.ScheduleDate,
		&currentTask.Deadline,
		&currentTask.Done,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}

	return currentTask, nil
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
