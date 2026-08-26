package task

type Repository interface {
	Add(title string) (Task, error)
	GetAll() ([]Task, error)
	GetByID(id int) (Task, error)
	MarkDone(id int) error
	Delete(id int) error
}
type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title string) (Task, error) {
	return s.repo.Add(title)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetTaskById(id int) (Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) DeleteTaskById(id int) error {
	return s.repo.Delete(id)
}

func (s *TaskService) MarkDoneTaskById(id int) error {
	return s.repo.MarkDone(id)
}
