package domain

import "time"

type Video struct {
	Id         string
	ResourceId string
	FilePath   string
	CreatedAt  time.Time
}
