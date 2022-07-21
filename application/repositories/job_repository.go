package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

type JobRepositoryInterface interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepository struct {
	Db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{Db: db}
}

func (repo JobRepository) Insert(job *domain.Job) (*domain.Job, error) {
	err := repo.Db.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repo JobRepository) Find(id string) (*domain.Job, error) {
	var job domain.Job

	repo.Db.Preload("Video").First(&job, "id = ?", id)

	if job.Id == "" {
		return nil, fmt.Errorf("job does not exists")
	}

	return &job, nil
}

func (repo JobRepository) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.Db.Save(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
