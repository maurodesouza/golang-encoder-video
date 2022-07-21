package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepositoryInterface interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepository struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{Db: db}
}

func (repo VideoRepository) Insert(video *domain.Video) (*domain.Video, error) {
	if video.Id == "" {
		video.Id = uuid.NewV4().String()
	}

	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (repo VideoRepository) Find(id string) (*domain.Video, error) {
	var video domain.Video

	repo.Db.First(&video, "id = ?", id)

	if video.Id == "" {
		return nil, fmt.Errorf("video does not exists")
	}

	return &video, nil
}
