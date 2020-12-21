package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DownloadFiles ..
func DownloadFiles(c *gin.Context) {

	// UserID := c.Param("user_id")
	// FolderID := c.Param("folder_id")
	FileID := c.Param("file_id")
	FilePath := "/media/main_transfer/" + FileID

	FileStat, err := os.Stat(FilePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	CDHeaferValue := fmt.Sprintf("attachment; filename=%s", FileID)
	FileSizeAsString := fmt.Sprintf("%d", FileStat.Size())

	c.Writer.Header().Add("Content-Disposition", CDHeaferValue)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Length", FileSizeAsString)
	c.File(FilePath)
}
