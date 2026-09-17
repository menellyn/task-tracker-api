package task

import "time"

type FakeRepository struct {
	tasks  []Task
	nextID int
	err    error
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		tasks:  make([]Task, 0, 5),
		nextID: 1,
	}
}

func (f *FakeRepository) Add(task Task) (Task, error) {
	if f.err != nil {
		return Task{}, f.err
	}

	newTask := Task{
		ID:           f.nextID,
		Title:        task.Title,
		Description:  task.Description,
		ScheduleDate: task.ScheduleDate,
		Deadline:     task.Deadline,
		Done:         false,
	}

	f.tasks = append(f.tasks, newTask)
	f.nextID++
	return newTask, nil
}

func (f *FakeRepository) GetAll() ([]Task, error) {
	if f.err != nil {
		return nil, f.err
	}
	returnedTasks := make([]Task, len(f.tasks))
	copy(returnedTasks, f.tasks)
	return returnedTasks, nil
}

func (f *FakeRepository) GetActual() ([]Task, error) {
	if f.err != nil {
		return nil, f.err
	}
	returnedTasks := make([]Task, 0, len(f.tasks))
	for _, task := range f.tasks {
		if !task.Done {
			returnedTasks = append(returnedTasks, task)
		}
	}
	return returnedTasks, nil
}

func (f *FakeRepository) GetByID(id int) (Task, error) {
	if f.err != nil {
		return Task{}, f.err
	}
	for _, task := range f.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}

func (f *FakeRepository) MarkDone(id int) error {
	if f.err != nil {
		return f.err
	}
	for i, task := range f.tasks {
		if task.ID == id {
			f.tasks[i].Done = true
			return nil
		}
	}
	return ErrTaskNotFound
}

func (f *FakeRepository) Update(id int, updateData map[string]interface{}) (Task, error) {
	if f.err != nil {
		return Task{}, f.err
	}
	currentTask := Task{}
	for i, task := range f.tasks {
		if task.ID == id {
			currentTask = f.tasks[i]
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
			f.tasks[i] = currentTask
			return currentTask, nil
		}
	}
	return Task{}, ErrTaskNotFound

}

func (f *FakeRepository) Delete(id int) error {
	if f.err != nil {
		return f.err
	}
	for i, task := range f.tasks {
		if task.ID == id {
			f.tasks = append(f.tasks[:i], f.tasks[i+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}
