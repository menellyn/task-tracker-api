package task

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func parseDate(value string) *time.Time {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return &date
}

func TestTaskService_CreateTask(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	taskToCreate := CreateTaskRequest{
		Title: "Buy milk",
	}
	got, err := s.CreateTask(taskToCreate)
	if err != nil {
		t.Fatal(err)
	}
	want := TaskResponse{
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
	taskToCreate := CreateTaskRequest{
		Title: " Buy milk  ",
	}
	got, err := s.CreateTask(taskToCreate)
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
	taskToCreate := CreateTaskRequest{
		Title: "Buy milk",
	}
	_, err := s.CreateTask(taskToCreate)
	if !errors.Is(err, repo.err) {
		t.Fatalf("got %v, expected %v", err, repo.err)
	}

}

func TestTaskService_CreateEmptyTask(t *testing.T) {
	repo := NewFakeRepository()
	s := NewTaskService(repo)
	tasksToCreate := []CreateTaskRequest{
		{
			Title: "",
		},
		{
			Title: "  ",
		},
	}
	for _, taskToCreate := range tasksToCreate {
		_, err := s.CreateTask(taskToCreate)
		if !errors.Is(err, ErrEmptyTitle) {
			t.Fatalf("got %v, expected %v", err, ErrEmptyTitle)
		}
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	repo := NewFakeRepository()
	repo.tasks = []Task{
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
	s := NewTaskService(repo)

	want := []TaskResponse{
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
	want := TaskResponse{
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
	want := TaskResponse{
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

func TestTaskService_UpdateTask(t *testing.T) {
	description := "Coconut milk"
	updatedTitle := "Buy milk!"
	updatedEmptyTitle := ""
	updatedDone := true
	cases := []struct {
		name          string
		id            int
		initTasks     []Task
		updateRequest UpdateTaskRequest
		want          TaskResponse
		err           error
	}{
		{
			name: "title",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Title: &updatedTitle,
			},
			want: TaskResponse{
				ID:    1,
				Title: "Buy milk!",
				Done:  false,
			},
			err: nil,
		},
		{
			name: "done",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Done: &updatedDone,
			},
			want: TaskResponse{
				ID:    1,
				Title: "Buy milk",
				Done:  true,
			},
			err: nil,
		},
		{
			name: "description",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Description: &description,
			},
			want: TaskResponse{
				ID:          1,
				Title:       "Buy milk",
				Description: &description,
				Done:        false,
			},
			err: nil,
		},
		{
			name: "deadline after schedule",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				ScheduleDate: parseDate("2026-09-20"),
				Deadline:     parseDate("2026-09-25"),
			},
			want: TaskResponse{
				ID:           1,
				Title:        "Buy milk",
				ScheduleDate: parseDate("2026-09-20"),
				Deadline:     parseDate("2026-09-25"),
				Done:         false,
			},
			err: nil,
		},
		{
			name: "deadline before schedule",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				ScheduleDate: parseDate("2026-09-20"),
				Deadline:     parseDate("2026-09-19"),
			},
			err: ErrInvalidTaskDates,
		},
		{
			name: "only deadline correct",
			id:   1,
			initTasks: []Task{
				{
					ID:           1,
					Title:        "Buy milk",
					ScheduleDate: parseDate("2026-09-20"),
					Done:         false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Deadline: parseDate("2026-09-25"),
			},
			want: TaskResponse{
				ID:           1,
				Title:        "Buy milk",
				ScheduleDate: parseDate("2026-09-20"),
				Deadline:     parseDate("2026-09-25"),
				Done:         false,
			},
			err: nil,
		},
		{
			name: "only schedule correct",
			id:   1,
			initTasks: []Task{
				{
					ID:       1,
					Title:    "Buy milk",
					Deadline: parseDate("2026-09-25"),
					Done:     false,
				},
			},
			updateRequest: UpdateTaskRequest{
				ScheduleDate: parseDate("2026-09-20"),
			},
			want: TaskResponse{
				ID:           1,
				Title:        "Buy milk",
				ScheduleDate: parseDate("2026-09-20"),
				Deadline:     parseDate("2026-09-25"),
				Done:         false,
			},
			err: nil,
		},
		{
			name: "invalid deadline",
			id:   1,
			initTasks: []Task{
				{
					ID:           1,
					Title:        "Buy milk",
					ScheduleDate: parseDate("2026-09-20"),
					Done:         false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Deadline: parseDate("2026-09-19"),
			},
			err: ErrInvalidTaskDates,
		},
		{
			name: "empty title",
			id:   1,
			initTasks: []Task{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			updateRequest: UpdateTaskRequest{
				Title: &updatedEmptyTitle,
			},
			err: ErrEmptyTitle,
		},
		{
			name:      "task not found",
			id:        1,
			initTasks: []Task{},
			updateRequest: UpdateTaskRequest{
				Description: &description,
			},
			err: ErrTaskNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewFakeRepository()
			repo.tasks = tc.initTasks
			service := NewTaskService(repo)
			got, err := service.UpdateTask(tc.id, tc.updateRequest)
			if !errors.Is(err, tc.err) {
				t.Errorf("Update() error = %v, wantErr %v", err, tc.err)
			}
			if tc.err == nil {
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("Update() got = %+v, want %+v", got, tc.want)
				}
			}
		})
	}

}
