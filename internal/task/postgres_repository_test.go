package task

import (
	"errors"
	"reflect"
	"testing"
)

func TestPostgresRepository_Add(t *testing.T) {
	tx := newTestTx(t)
	repo := NewPostgresRepository(tx)
	got, err := repo.Add("Buy milk")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID <= 0 {
		t.Fatalf("generated ID should be > 0, got %d", got.ID)
	}
	wantTitle := "Buy milk"
	if got.Title != wantTitle {
		t.Fatalf("got title %q, want %q", got.Title, wantTitle)
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
	if title != wantTitle {
		t.Fatalf("got title %q in database, want %q", title, wantTitle)
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
			createdTask, err := repo.Add("Buy milk")
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

	wantTitles := []string{"Buy milk", "Drink coffee", "Wash dishes"}

	for _, title := range wantTitles {
		_, err := repo.Add(title)
		if err != nil {
			t.Fatal(err)
		}
	}

	tasks, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != len(wantTitles) {
		t.Fatalf("Want %d tasks, got %d tasks", len(wantTitles), len(tasks))
	}

	for i := range wantTitles {
		if wantTitles[i] != tasks[i].Title {
			t.Fatalf("%d tasks: want %q, got %q", i, wantTitles[i], tasks[i].Title)
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
	task, err := repo.Add("Buy milk")
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

	task, err := repo.Add("Buy milk")
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

	titles := []string{"Buy milk", "Wash dishes", "Read book"}
	tasks := []Task{}

	for _, title := range titles {
		task, err := repo.Add(title)
		if err != nil {
			t.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	wantTitles := []string{"Buy milk", "Read book"}
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

	if len(got) != len(wantTitles) {
		t.Fatalf("expected %d tasks, got %d tasks", len(wantTitles), len(got))
	}

	for i := range wantTitles {
		if wantTitles[i] != got[i].Title {
			t.Fatalf("%d task: want %q, got %q", i, wantTitles[i], got[i].Title)
		}
	}
}
