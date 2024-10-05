/*
* Editor:KaiWen
* PATH:/main.go
* Description:
 */

package main

import (
	"gardeManageSystem/handlers"
	"gardeManageSystem/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	models.StudentMap = make(map[string]*models.Student)

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// 路由配置
	r.POST("/api/upload", handlers.UploadExcel)
	r.GET("/api/students", handlers.GetStudents)
	r.GET("/api/analysis", handlers.GetAnalysis)
	r.GET("/api/upload-history", handlers.GetUploadHistory)
	r.GET("/api/student-analysis/:id", handlers.GetStudentAnalysis)

	log.Fatal(r.Run(":8080"))
}
