package app

import (
	"github.com/petrakypetrov/main_transfer_api/controllers"
)

func mapUrls() {

	router.GET("/ping", controllers.Ping)

	// group: v1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users/:user_id/folders/:folder_id/files/:file_id", controllers.DownloadFiles)
		v1.GET("/users/:user_id/folders/:folder_id/files", controllers.DownloadFiles)
		v1.POST("/users/:user_id/folders/:folder_id/files", controllers.UploadFiles)
	}
}
