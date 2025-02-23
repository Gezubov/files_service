package service

import (
	"github.com/Gezubov/file_service/internal/models"
	"github.com/google/uuid"
)

type FileRepository interface {
	Create(file *models.File) error
	GetByUUID(uuid uuid.UUID) (*models.File, error)
	Delete(uuid uuid.UUID) error
	GetAllFiles() ([]models.File, error)
}
