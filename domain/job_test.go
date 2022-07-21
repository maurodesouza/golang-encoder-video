package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()

	video.Id = uuid.NewV4().String()
	video.ResourceId = "abc"
	video.FilePath = "abc"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("test", "test", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
