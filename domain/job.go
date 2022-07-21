package domain

import "time"

type Job struct {
	Id               string
	OutputBucketPath string
	Status           string
	Video            *Video
	VideoId          string
	Error            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
