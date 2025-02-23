package service

import (
	"errors"

	"github.com/Gezubov/file_service/internal/models"

	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type FileService struct {
	fileRepo FileRepository
}

func NewFileService(fileRepo FileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}

func (s *FileService) CreateFile(file *models.File) error {
	existingFileByUUID, err := s.fileRepo.GetByUUID(file.UUID)
	if err == nil && existingFileByUUID != nil {
		return errors.New("UUID уже занят")
	}

	return s.fileRepo.Create(file)
}

func (s *FileService) GetFileByID(uuid uuid.UUID) (*models.File, error) {
	return s.fileRepo.GetByUUID(uuid)
}

func (s *FileService) DeleteFile(uuid uuid.UUID) error {
	return s.fileRepo.Delete(uuid)
}

func (s *FileService) GetAllFiles() ([]models.File, error) {
	return s.fileRepo.GetAllFiles()
}
