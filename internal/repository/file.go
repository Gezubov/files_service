package repository

import (
	"database/sql"

	"log/slog"
	"time"

	"github.com/Gezubov/file_service/internal/models"
	"github.com/google/uuid"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {
	slog.Info("Creating file", "file_name", file.FileName)
	query := `
		INSERT INTO files (uuid, file_name, file_size, file_link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid`

	now := time.Now()
	file.UUID = uuid.New()
	err := r.db.QueryRow(
		query,
		file.UUID,
		file.FileName,
		file.FileSize,
		file.FileLink,
		now,
		now,
	).Scan(&file.UUID)

	if err != nil {
		slog.Error("Error creating user", "error", err)
		return err
	}

	file.CreatedAt = now
	file.UpdatedAt = now
	return nil
}

func (r *FileRepository) GetByUUID(uuid uuid.UUID) (*models.File, error) {
	slog.Info("Getting file with UUID", "uuid", uuid)
	file := &models.File{}

	query := `
		SELECT uuid, file_name, file_size, file_link, created_at, updated_at
		FROM files
		WHERE uuid = $1`

	err := r.db.QueryRow(query, uuid).Scan(
		&file.UUID,
		&file.FileName,
		&file.FileSize,
		&file.FileLink,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		slog.Warn("File not found", "uuid", uuid)
		return nil, ErrFileNotFound
	}
	if err != nil {
		slog.Error("Error fetching file", "uuid", uuid, "error", err)
		return nil, err
	}

	return file, nil
}

func (r *FileRepository) Delete(uuid uuid.UUID) error {
	slog.Info("Deleting file", "uuid", uuid)
	query := `DELETE FROM files WHERE uuid = $1`

	result, err := r.db.Exec(query, uuid)
	if err != nil {
		slog.Error("Error deleting file", "uuid", uuid, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error checking delete result", "uuid", uuid, "error", err)
		return err
	}
	if rowsAffected == 0 {
		slog.Warn("File not found during deletion", "uuid", uuid)
		return ErrFileNotFound
	}

	return nil
}

func (r *FileRepository) GetAllFiles() ([]models.File, error) {
	slog.Info("Fetching all files from database")

	query := `SELECT uuid, file_name, file_size, file_link, created_at, updated_at FROM files`
	rows, err := r.db.Query(query)

	if err != nil {
		slog.Error("Error executing query to fetch files", "error", err)
		return nil, err
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.UUID, &file.FileName, &file.FileSize, &file.FileLink, &file.CreatedAt, &file.UpdatedAt); err != nil {
			slog.Error("Error scanning file row", "error", err)
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over user rows", "error", err)
		return nil, err
	}

	slog.Info("Successfully fetched files", "count", len(files))
	return files, nil
}
