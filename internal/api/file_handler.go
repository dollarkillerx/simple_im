package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simple_im/internal/middleware"
	"simple_im/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (a *ApiServer) Upload(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Check file size
	if file.Size > a.conf.UploadConfiguration.MaxSize {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("file size exceeds limit (%d bytes)", a.conf.UploadConfiguration.MaxSize),
		})
		return
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	allowed := false
	for _, t := range a.conf.UploadConfiguration.AllowTypes {
		if t == contentType {
			allowed = true
			break
		}
	}
	if !allowed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	// Create upload directory if not exists
	uploadDir := a.conf.UploadConfiguration.SavePath
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Error().Err(err).Msg("failed to create upload directory")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, newFilename)

	// Save file
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		log.Error().Err(err).Msg("failed to save uploaded file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Save file record to database
	fileRecord := &models.File{
		UserID:   userID,
		Filename: file.Filename,
		Filepath: savePath,
		Filesize: file.Size,
		Mimetype: contentType,
	}

	if err := a.storage.GetDB().Create(fileRecord).Error; err != nil {
		log.Error().Err(err).Msg("failed to save file record")
		// File saved but record failed, still return success
	}

	// Generate URL
	fileURL := fmt.Sprintf("/files/%s", newFilename)

	ctx.JSON(http.StatusOK, gin.H{
		"id":       fileRecord.ID,
		"filename": file.Filename,
		"url":      fileURL,
		"size":     file.Size,
		"mimetype": contentType,
	})
}

// GetFileType returns the message type based on content type
func GetFileType(contentType string) models.MessageType {
	if strings.HasPrefix(contentType, "image/") {
		return models.MsgTypeImage
	}
	return models.MsgTypeFile
}
