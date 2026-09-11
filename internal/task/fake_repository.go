package task

type fakeRepository struct {
	tasks  []Task
	nextID int
	err    error
}

func NewFakeRepository() *fakeRepository {
	return &fakeRepository{
		tasks:  make([]Task, 0, 5),
		nextID: 1,
	}
}

func (f *fakeRepository) Add(title string) (Task, error) {
	if f.err != nil {
		return Task{}, f.err
	}

	newTask := Task{
		ID:    f.nextID,
		Title: title,
		Done:  false,
	}
	f.tasks = append(f.tasks, newTask)
	f.nextID++
	return newTask, nil
}

func (f *fakeRepository) GetAll() ([]Task, error) {
	if f.err != nil {
		return nil, f.err
	}
	returnedTasks := make([]Task, len(f.tasks))
	copy(returnedTasks, f.tasks)
	return returnedTasks, nil
}

func (f *fakeRepository) GetActual() ([]Task, error) {
	if f.err != nil {
		return nil, f.err
	}
	returnedTasks := make([]Task, len(f.tasks))
	for _, task := range f.tasks {
		if !task.Done {
			returnedTasks = append(returnedTasks, task)
		}
	}
	return returnedTasks, nil
}

func (f *fakeRepository) GetByID(id int) (Task, error) {
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

func (f *fakeRepository) MarkDone(id int) error {
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

func (f *fakeRepository) Update(taskToUpdate Task) (Task, error) {
	if f.err != nil {
		return Task{}, f.err
	}
	for i, task := range f.tasks {
		if task.ID == taskToUpdate.ID {
			f.tasks[i] = taskToUpdate
			return taskToUpdate, nil
		}
	}
	return Task{}, ErrTaskNotFound

}

func (f *fakeRepository) Delete(id int) error {
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
