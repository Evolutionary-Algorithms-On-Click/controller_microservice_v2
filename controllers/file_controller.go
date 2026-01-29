package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type FileController struct {
	Module *modules.FileModule
	Logger zerolog.Logger
}

func NewFileController(module *modules.FileModule, logger zerolog.Logger) *FileController {
	return &FileController{
		Module: module,
		Logger: logger,
	}
}

func (c *FileController) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	//CHANGEABLE: Limit upload size 50MB 
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path, err := c.Module.UploadFile(sessionID, file, header)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to upload file")
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusCreated, map[string]string{
		"path":     path,
		"filename": filepath.Base(path),
	}, &c.Logger)
}

func (c *FileController) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	files, err := c.Module.ListFiles(sessionID)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to list files")
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSONResponseWithLogger(w, http.StatusOK, map[string][]string{
		"files": files,
	}, &c.Logger)
}
