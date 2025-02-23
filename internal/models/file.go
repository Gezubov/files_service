package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	UUID      uuid.UUID `json:"uuid"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	FileLink  string    `json:"file_link"` //link to s3 bucket
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
