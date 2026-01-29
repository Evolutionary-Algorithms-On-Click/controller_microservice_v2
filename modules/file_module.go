package modules

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileModule struct {
	BaseDir string
}

func NewFileModule(baseDir string) *FileModule {
	return &FileModule{
		BaseDir: baseDir,
	}
}

func (m *FileModule) UploadFile(sessionID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error) {
	sessionDir := filepath.Join(m.BaseDir, sessionID.String())
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create session directory: %w", err)
	}

	// Sanitize filename to prevent directory traversal
	filename := filepath.Base(header.Filename)
	filename = strings.ReplaceAll(filename, "..", "") // simple extra safety

	filePath := filepath.Join(sessionDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	return filePath, nil
}

func (m *FileModule) ListFiles(sessionID uuid.UUID) ([]string, error) {
	sessionDir := filepath.Join(m.BaseDir, sessionID.String())
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (m *FileModule) DeleteSessionFiles(sessionID uuid.UUID) error {
	sessionDir := filepath.Join(m.BaseDir, sessionID.String())
	// Check if exists first to avoid error if already gone
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(sessionDir); err != nil {
		return fmt.Errorf("failed to delete session directory: %w", err)
	}
	return nil
}
