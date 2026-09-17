package task

import (
	"errors"

	"gorm.io/gorm"
)

type ORMRepository struct {
	db *gorm.DB
}

func NewORMRepository(db *gorm.DB) *ORMRepository {
	return &ORMRepository{
		db: db,
	}
}

func (r *ORMRepository) Add(task Task) (Task, error) {
	if err := r.db.Create(&task).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *ORMRepository) GetByID(id int) (Task, error) {
	task := Task{}

	err := r.db.Where("id = ?", id).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Task{}, ErrTaskNotFound
	}

	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *ORMRepository) GetAll() ([]Task, error) {
	var tasks []Task
	if err := r.db.Order("id").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ORMRepository) GetActual() ([]Task, error) {
	var tasks []Task
	if err := r.db.Where("done = ?", false).Order("id").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *ORMRepository) Update(id int, updateData map[string]interface{}) (Task, error) {
	result := r.db.Model(&Task{}).Where("id = ?", id).Updates(updateData)

	if result.Error != nil {
		return Task{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Task{}, ErrTaskNotFound
	}

	return r.GetByID(id)
}

func (r *ORMRepository) MarkDone(id int) error {
	result := r.db.Model(&Task{}).Where("id = ?", id).Update("done", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *ORMRepository) Delete(id int) error {
	result := r.db.Delete(&Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
