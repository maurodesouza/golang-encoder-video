package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDbTest()

	defer db.Close()

	video := domain.NewVideo()

	video.Id = uuid.NewV4().String()
	video.FilePath = "video.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepository{Db: db}

	return video, repo
}

func TestVideoServiceDownloadAndFragment(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()

	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("microservice_encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}
