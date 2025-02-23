package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Gezubov/file_service/internal/models"
	"github.com/Gezubov/file_service/internal/service"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type FileController struct {
	fileService *service.FileService
}

func NewFileController(fileService *service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

func (c *FileController) CreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var file models.File
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.fileService.CreateFile(&file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *FileController) GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := c.fileService.GetFileByID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(file)
}

func (c *FileController) GetFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := c.fileService.GetAllFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}

func (c *FileController) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid File ID", http.StatusBadRequest)
		return
	}

	if err := c.fileService.DeleteFile(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *FileController) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем файл из запроса
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Здесь вы можете сохранить файл на диск или сразу загрузить в MinIO
	// Пример загрузки в MinIO
	minioService, err := service.NewMinioService("minio:9000", "minioadmin", "minioadmin", "your-bucket-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	objectName := "uploaded-file" // Замените на уникальное имя файла
	if err := minioService.UploadFile(file, objectName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
