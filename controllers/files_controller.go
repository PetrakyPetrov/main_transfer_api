package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/petrakypetrov/main_transfer_api/config"
)

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// DownloadFiles ..
func DownloadFiles(c *gin.Context) {

	// UserID := c.Param("user_id")
	// FolderID := c.Param("folder_id")
	FileID := c.Param("file_id")
	FilePath := config.WorkDir + FileID

	if FileID == "" {
		ZipName := randomString(10) + ".zip"
		c.Writer.Header().Set("Content-type", "application/octet-stream")
		c.Stream(func(w io.Writer) bool {
			// Create a zip archive.
			ar := zip.NewWriter(w)
			walker := func(path string, info os.FileInfo, err error) error {
				// fmt.Printf("Crawling: %#v\n", path)
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				fName := info.Name()
				fExt := filepath.Ext(fName)
				if fExt == ".zip" {
					return nil
				}

				f, err := ar.Create(fName)
				if err != nil {
					return err
				}

				c.Writer.Header().Set("Content-Disposition", "attachment; filename="+ZipName)
				_, err = io.Copy(f, file)
				if err != nil {
					return err
				}

				return nil
			}

			err := filepath.Walk(config.WorkDir, walker)
			if err != nil {
				panic(err)
			}

			ar.Close()
			return false
		})
		return
	}

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

// UploadFiles ...
func UploadFiles(c *gin.Context) {

	form, _ := c.MultipartForm()
	files := form.File["files"]

	// for _, file := range files {
	// 	log.Println(file.Filename)

	// 	// Upload the file to specific dst.
	// 	c.SaveUploadedFile(file, config.WorkDir)
	// }

	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Println("error opening file ", err)
		}

		dst, err := os.Create(config.WorkDir + files[i].Filename)
		defer dst.Close()
		if err != nil {
			fmt.Println("error creating destination ", err)
		}

		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println("error copying file", err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
	return
}
