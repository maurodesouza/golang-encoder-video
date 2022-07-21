package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryInsert(t *testing.T) {
	db := database.NewDbTest()

	video := domain.NewVideo()

	video.Id = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)

	repo.Insert(video)
	v, err := repo.Find(video.Id)

	require.Nil(t, err)
	require.NotEmpty(t, v)
	require.Equal(t, v.Id, video.Id)
}
