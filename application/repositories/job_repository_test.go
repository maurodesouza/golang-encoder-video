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

func TestJobRepositoryInsert(t *testing.T) {
	db := database.NewDbTest()

	video := domain.NewVideo()

	video.Id = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repoVideo := repositories.NewVideoRepository(db)
	repoVideo.Insert(video)

	job, err := domain.NewJob("path", "testing", video)

	require.Nil(t, err)
	require.NotEmpty(t, job.Id)

	repoJob := repositories.NewJobRepository(db)
	repoJob.Insert(job)

	j, err := repoJob.Find(job.Id)

	require.Nil(t, err)
	require.NotEmpty(t, j)

	require.Equal(t, j.VideoId, video.Id)
	require.Equal(t, j.Id, job.Id)
}

func TestJobRepositoryUpdate(t *testing.T) {
	db := database.NewDbTest()

	video := domain.NewVideo()

	video.Id = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repoVideo := repositories.NewVideoRepository(db)
	repoVideo.Insert(video)

	job, err := domain.NewJob("path", "testing", video)

	require.Nil(t, err)
	require.NotEmpty(t, job.Id)

	repoJob := repositories.NewJobRepository(db)
	repoJob.Insert(job)

	newStatus, newOutputPath := "new status", "new output path"

	job.Status = newStatus
	job.OutputBucketPath = newOutputPath

	j, err := repoJob.Update(job)

	require.Nil(t, err)
	require.NotEmpty(t, j)

	require.Equal(t, j.Status, newStatus)
	require.Equal(t, j.OutputBucketPath, newOutputPath)
}
