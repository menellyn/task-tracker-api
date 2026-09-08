package task

import (
	"errors"
	"reflect"
	"testing"
)

func TestORMRepository_Add(t *testing.T) {
	tx := newTestGormTx(t)
	repo := NewORMRepository(tx)

	wantTitle := "Buy milk"

	got, err := repo.Add(wantTitle)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID <= 0 {
		t.Errorf("got %d, want a positive ID", got.ID)
	}

	if got.Title != wantTitle {
		t.Errorf("got title %s, want %s", got.Title, wantTitle)
	}

	if got.Done {
		t.Error("new task should be not done")
	}

}

func TestORMRepository_GetByID(t *testing.T) {
	cases := []struct {
		name string
		id   int
		want Task
		err  error
	}{
		{
			name: "task exists",
			want: Task{
				Title: "Buy milk",
				Done:  false,
			},
			err: nil,
		},
		{
			name: "task not found",
			id:   1,
			err:  ErrTaskNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := newTestGormTx(t)
			repo := NewORMRepository(tx)

			if tc.err == nil {
				createdTask, err := repo.Add(tc.want.Title)
				if err != nil {
					t.Fatal(err)
				}
				tc.want.ID = createdTask.ID
				tc.id = createdTask.ID
			}

			got, err := repo.GetByID(tc.id)
			if !errors.Is(err, tc.err) {
				t.Fatalf("got error %v, want %v", tc.err, err)
			}

			if tc.err == nil {
				if tc.want != got {
					t.Errorf("want %+v, got %+v", tc.want, got)
				}
			}
		})
	}

}

func TestORMRepository_GetAll(t *testing.T) {
	cases := []struct {
		name   string
		titles []string
		want   []Task
	}{
		{
			name: "all tasks",
			titles: []string{
				"Buy milk",
				"Wash dishes",
				"Evening walk",
			},
			want: []Task{
				{
					Title: "Buy milk",
					Done:  false,
				},
				{
					Title: "Wash dishes",
					Done:  false,
				},
				{
					Title: "Evening walk",
					Done:  false,
				},
			},
		},
		{
			name: "no tasks",
			want: []Task{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := newTestGormTx(t)
			repo := NewORMRepository(tx)

			if tc.titles != nil {
				for i, title := range tc.titles {
					createdTask, err := repo.Add(title)
					if err != nil {
						t.Fatal(err)
					}
					tc.want[i].ID = createdTask.ID
				}
			}

			got, err := repo.GetAll()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestORMRepository_MarkDone(t *testing.T) {
	cases := []struct {
		name       string
		id         int
		want       Task
		taskExists bool
		err        error
	}{
		{
			name: "task exists",
			want: Task{
				Title: "Buy milk",
				Done:  true,
			},
			taskExists: true,
			err:        nil,
		},
		{
			name:       "task not found",
			id:         1,
			taskExists: false,
			err:        ErrTaskNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := newTestGormTx(t)
			repo := NewORMRepository(tx)
			if tc.taskExists {
				createdTask, err := repo.Add(tc.want.Title)
				if err != nil {
					t.Fatal(err)
				}
				tc.want.ID = createdTask.ID
				tc.id = createdTask.ID
			}

			err := repo.MarkDone(tc.id)
			if !errors.Is(err, tc.err) {
				t.Fatalf("got error %v, want %v", tc.err, err)
			}

			if tc.taskExists {
				got, err := repo.GetByID(tc.id)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %+v, want %+v", got, tc.want)
				}
			}
		})
	}
}

func TestORMRepository_Delete(t *testing.T) {
	cases := []struct {
		name       string
		id         int
		deletedIdx int
		titles     []string
		want       []Task
		taskExists bool
		err        error
	}{
		{
			name: "task exists",
			titles: []string{
				"Buy milk",
			},
			deletedIdx: 0,
			want:       []Task{},
			taskExists: true,
			err:        nil,
		},
		{
			name:       "task not found",
			id:         9999,
			taskExists: false,
			err:        ErrTaskNotFound,
		},
		{
			name: "task delete doesn't delete other task",
			titles: []string{
				"Buy milk",
				"Wash dishes",
				"Evening walk",
			},
			deletedIdx: 1,
			want: []Task{
				{
					Title: "Buy milk",
					Done:  false,
				},
				{
					Title: "Evening walk",
					Done:  false,
				},
			},
			taskExists: true,
			err:        nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := newTestGormTx(t)
			repo := NewORMRepository(tx)
			if tc.taskExists {
				createdTasks := make([]Task, 0, len(tc.titles))
				for _, title := range tc.titles {
					createdTask, err := repo.Add(title)
					if err != nil {
						t.Fatal(err)
					}
					createdTasks = append(createdTasks, createdTask)
				}
				tc.id = createdTasks[tc.deletedIdx].ID
			}

			err := repo.Delete(tc.id)
			if !errors.Is(err, tc.err) {
				t.Fatalf("got error %v, want %v", tc.err, err)
			}

			if tc.taskExists {
				got, err := repo.GetAll()
				if err != nil {
					t.Fatal(err)
				}
				if len(got) != len(tc.want) {
					t.Fatalf("got %d tasks, want %d", len(got), len(tc.want))
				}

				if len(got) != 0 {
					for i, gotTask := range got {
						wantTask := tc.want[i]
						if gotTask.Title != wantTask.Title {
							t.Errorf("got title %s, want %s", gotTask.Title, wantTask.Title)
						}
						if gotTask.Done != wantTask.Done {
							t.Errorf("got done %v, want %v", gotTask.Done, wantTask.Done)
						}
					}
				}
			}
		})
	}
}
