package task

import (
	"errors"
	"reflect"
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	got, err := s.CreateTask("Buy milk")
	if err != nil {
		t.Fatal(err)
	}
	want := Task{
		Title: "Buy milk",
		Done:  false,
	}
	if got.ID <= 0 {
		t.Fatalf("got ID %d, want positive ID", got.ID)
	}
	if got.Title != want.Title {
		t.Fatalf("got %v, want %v", got.Title, want.Title)
	}
	if got.Done != want.Done {
		t.Fatalf("got %v, want %v", got.Done, want.Done)
	}
}

func TestTaskService_CreateTaskTrimTitle(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	got, err := s.CreateTask(" Buy milk  ")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Buy milk" {
		t.Fatalf("got %q, want %q", got.Title, "Buy milk")
	}
}

func TestTaskService_CreateTaskWithError(t *testing.T) {
	repo := NewFakeRepository()
	repo.err = errors.New("database error")
	s := NewTaskService(repo)
	_, err := s.CreateTask("Buy milk")
	if !errors.Is(err, repo.err) {
		t.Fatalf("got %v, expected %v", err, repo.err)
	}

}

func TestTaskService_CreateEmptyTask(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	_, err := s.CreateTask("")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("got %v, expected %v", err, ErrEmptyTitle)
	}
	_, err = s.CreateTask("  ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("got %v, expected %v", err, ErrEmptyTitle)
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	repo := NewFakeRepository()
	want := []Task{
		{
			ID:    1,
			Title: "Buy milk",
			Done:  false,
		},
		{
			ID:    2,
			Title: "Drink hot chocolate",
			Done:  false,
		},
		{
			ID:    3,
			Title: "Wash dishes",
			Done:  false,
		},
	}
	repo.tasks = want
	s := NewTaskService(repo)

	got, err := s.GetAllTasks()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestTaskService_GetAllTasksWithError(t *testing.T) {
	repo := NewFakeRepository()
	repo.err = errors.New("database error")
	s := NewTaskService(repo)

	_, err := s.GetAllTasks()
	if !errors.Is(err, repo.err) {
		t.Fatalf("got %v, expected %v", err, repo.err)
	}
}

func TestTaskService_GetByID(t *testing.T) {
	repo := NewFakeRepository()
	repo.tasks = []Task{
		{
			ID:    1,
			Title: "Buy milk",
			Done:  false,
		},
	}
	s := NewTaskService(repo)
	want := Task{
		ID:    1,
		Title: "Buy milk",
		Done:  false,
	}
	got, err := s.GetTaskByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTaskService_GetByIDNotFound(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	_, err := s.GetTaskByID(100)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("got %v, expected %v", err, ErrTaskNotFound)
	}
}

func TestTaskService_MarkDoneTaskByID(t *testing.T) {
	repo := NewFakeRepository()
	repo.tasks = []Task{
		{
			ID:    1,
			Title: "Buy milk",
			Done:  false,
		},
	}
	s := NewTaskService(repo)
	want := Task{
		ID:    1,
		Title: "Buy milk",
		Done:  true,
	}
	err := s.MarkDoneTaskByID(1)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.GetTaskByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

}

func TestTaskService_MarkDoneTaskByIDNotFound(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	err := s.MarkDoneTaskByID(100)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("got %v, expected %v", err, ErrTaskNotFound)
	}
}

func TestTaskService_DeleteTaskByID(t *testing.T) {
	repo := NewFakeRepository()
	repo.tasks = []Task{
		{
			ID:    1,
			Title: "Buy milk",
			Done:  false,
		},
	}
	s := NewTaskService(repo)

	if err := s.DeleteTaskByID(1); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetTaskByID(1)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("got %v, expected %v", err, ErrTaskNotFound)
	}

}

func TestTaskService_DeleteTaskByIDNotFound(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	err := s.DeleteTaskByID(100)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("got %v, expected %v", err, ErrTaskNotFound)
	}
}
