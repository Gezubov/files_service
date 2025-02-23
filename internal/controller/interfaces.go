package controller

import (
	"github.com/Gezubov/file_service/internal/models"
	"github.com/google/uuid"
)

type FileService interface {
	CreateFile(file *models.File) error
	GetFileByID(uuid uuid.UUID) (*models.File, error)
	DeleteFile(id int64) error
	GetAllFiles() ([]models.File, error)
}
