/*
* Editor:KaiWen
* PATH:handlers/grade_handler.go
* Description:
 */

package handlers

import (
	"fmt"
	"gardeManageSystem/models"
	"gardeManageSystem/services"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 全局变量
var (
	studentList []models.Student // 改名避免与类型名冲突
	mutexLock   sync.RWMutex     // 改名更具描述性
)

func UploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File upload failed"})
		return
	}
	// 保存文件到服务器
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "Failed to save file", "data": nil})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "Cannot open file"})
		return
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(400, gin.H{"error": "Cannot read Excel file"})
		return
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		c.JSON(400, gin.H{"error": "Cannot read sheet"})
		return
	}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 10 {
			continue
		}

		studentID := row[1]
		student, exists := models.StudentMap[studentID]
		if !exists {
			student = &models.Student{
				Name: row[0],
				ID:   studentID,
			}
			models.StudentMap[studentID] = student
		}

		grade := models.Grade{
			ExamName:  file.Filename, // 使用文件名作为考试名称
			Chinese:   services.ParseFloat(row[2]),
			Math:      services.ParseFloat(row[3]),
			English:   services.ParseFloat(row[4]),
			Chemistry: services.ParseFloat(row[5]),
			Physics:   services.ParseFloat(row[6]),
			Geography: services.ParseFloat(row[7]),
			Biology:   services.ParseFloat(row[8]),
			Total:     services.ParseFloat(row[9]),
		}
		student.Grades = append(student.Grades, grade)
	}

	// 保存学生数据到JSON文件
	services.SaveStudentsToFile()

	// 保存上传历史
	uploadHistory := models.UploadHistory{
		ID:         len(models.UploadHistoryList) + 1,
		FileName:   file.Filename,
		UploadTime: time.Now(),
		Status:     "ok",
	}
	models.UploadHistoryList = append(models.UploadHistoryList, uploadHistory)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Upload successful",
		"data":    uploadHistory,
	})
}

func GetStudents(c *gin.Context) {
	// 从JSON文件中读取学生数据
	services.LoadStudentsFromFile()
	fmt.Println(models.StudentMap)
	c.JSON(200, gin.H{
		"code": 200,
		"data": models.StudentMap,
	})
}

func GetAnalysis(c *gin.Context) {
	if len(models.StudentMap) == 0 {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "No data available",
			"data":    nil,
		})
		return
	}

	students := make([]*models.Student, 0, len(models.StudentMap))
	for _, student := range models.StudentMap {
		students = append(students, student)
	}

	analysis := services.AnalyzeGrades(students)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success",
		"data":    analysis,
	})
}

func GetUploadHistory(c *gin.Context) {
	// 清空列表
	models.UploadHistoryList = []models.UploadHistory{}

	// 获取uploads目录中的文件
	err := filepath.Walk("./uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 过滤文件夹，只处理文件
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".xlsx") || strings.HasSuffix(info.Name(), ".xls") || strings.HasSuffix(info.Name(), ".csv")) {
			models.UploadHistoryList = append(models.UploadHistoryList, models.UploadHistory{
				ID:         len(models.UploadHistoryList) + 1,
				FileName:   info.Name(),
				UploadTime: info.ModTime(),
				Status:     "ok",
			})
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "读取上传历史失败", "data": nil})
		return
	}

	c.JSON(200, gin.H{
		"code":  200,
		"files": models.UploadHistoryList,
	})
}

func GetStudentAnalysis(c *gin.Context) {
	studentID := c.Param("id")
	fmt.Println("id", studentID)

	// 从JSON文件中读取学生数据
	services.LoadStudentsFromFile()

	student, exists := models.StudentMap[studentID]
	fmt.Println(student)
	if !exists {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "Student not found",
			"data":    nil,
		})
		return
	}

	analysis := services.AnalyzeStudentGrades(student)
	fmt.Println(analysis)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success",
		"data":    analysis,
	})
}
