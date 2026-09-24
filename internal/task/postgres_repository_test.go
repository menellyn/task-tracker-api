package task

import (
	"errors"
	"reflect"
	"testing"
)

func TestPostgresRepository_Add(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)
	initTask := Task{
		Title: "Buy milk",
	}
	got, err := repo.Add(initTask)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID <= 0 {
		t.Fatalf("generated ID should be > 0, got %d", got.ID)
	}
	if got.Title != initTask.Title {
		t.Fatalf("got title %q, want %q", got.Title, initTask.Title)
	}
	if got.Done {
		t.Fatal("new task should be not done")
	}

	var title string

	err = tx.QueryRow(`
		SELECT title
		FROM tasks
		WHERE id = $1
	`, got.ID).Scan(&title)

	if err != nil {
		t.Fatal(err)
	}
	if title != initTask.Title {
		t.Fatalf("got title %q in database, want %q", title, initTask.Title)
	}

}

func TestPostgresRepository_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		wantErr  error
		wantTask Task
	}{
		{
			name: "task exist",
			wantTask: Task{
				Title: "Buy milk",
				Done:  false,
			},
		},
		{
			name:    "task not found",
			id:      100,
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := newTestTx(t)
			repo := NewPostgresRepository(tx)
			createdTask, err := repo.Add(tt.wantTask)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantErr == nil {
				tt.id = createdTask.ID
				tt.wantTask.ID = createdTask.ID
			}

			gotTask, err := repo.GetByID(tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got %v error, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(gotTask, tt.wantTask) {
				t.Fatalf("got %+v, want %+v", gotTask, tt.wantTask)
			}
		})
	}
}

func TestPostgresRepository_GetAll(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)

	wantTasks := []Task{
		{
			Title: "Buy milk",
			Done:  false,
		},
		{
			Title: "Drink coffee",
			Done:  false,
		},
		{
			Title: "Wash dishes",
			Done:  false,
		},
	}

	for _, task := range wantTasks {
		_, err := repo.Add(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	gotTasks, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(gotTasks) != len(wantTasks) {
		t.Fatalf("Want %d tasks, got %d tasks", len(wantTasks), len(gotTasks))
	}

	for i := range wantTasks {
		if wantTasks[i].Title != gotTasks[i].Title {
			t.Fatalf("%d tasks: want %q, got %q", i, wantTasks[i].Title, gotTasks[i].Title)
		}
	}

}

func TestPostgresRepository_GetAllEmpty(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)

	want := []Task{}
	got, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %+v, got %+v", want, got)
	}

}

func TestPostgresRepository_MarkDone(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)
	initTask := Task{
		Title: "Buy milk",
		Done:  false,
	}
	task, err := repo.Add(initTask)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.MarkDone(task.ID); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Done {
		t.Fatal("expected task to be done")
	}

}

func TestPostgresRepository_MarkDoneNotFound(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)

	if err := repo.MarkDone(100); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestPostgresRepository_Delete(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)
	initTask := Task{
		Title: "Buy milk",
		Done:  false,
	}
	task, err := repo.Add(initTask)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Delete(task.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := repo.GetByID(task.ID); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestPostgresRepository_DeleteNotFound(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)

	if err := repo.Delete(100); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestPostgresRepository_DeleteDoesNotDeleteOtherTasks(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)

	tasks := []Task{}
	initTasks := []Task{
		{
			Title: "Buy milk",
			Done:  false,
		},
		{
			Title: "Wash dishes",
			Done:  false,
		},
		{
			Title: "Read book",
			Done:  false,
		},
	}

	for _, initTask := range initTasks {
		task, err := repo.Add(initTask)
		if err != nil {
			t.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	wantTasks := []Task{
		{
			Title: "Buy milk",
			Done:  false,
		},
		{
			Title: "Read book",
			Done:  false,
		},
	}
	if err := repo.Delete(tasks[1].ID); err != nil {
		t.Fatal(err)
	}

	if _, err := repo.GetByID(tasks[1].ID); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}

	got, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(wantTasks) {
		t.Fatalf("expected %d tasks, got %d tasks", len(wantTasks), len(got))
	}

	for i := range wantTasks {
		if wantTasks[i].Title != got[i].Title {
			t.Fatalf("%d task: want %q, got %q", i, wantTasks[i].Title, got[i].Title)
		}
	}
}
