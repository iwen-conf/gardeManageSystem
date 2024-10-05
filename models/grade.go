/*
* Editor:KaiWen
* PATH:models/grade.go
* Description:
 */

package models

type Grade struct {
	Chinese   float64 `json:"chinese"`
	Math      float64 `json:"math"`
	English   float64 `json:"english"`
	Chemistry float64 `json:"chemistry"`
	Physics   float64 `json:"physics"`
	Geography float64 `json:"geography"`
	Biology   float64 `json:"biology"`
	Total     float64 `json:"total"`
	ExamName  string  `json:"examName"`
}
type Student struct {
	Name   string  `json:"name"`
	ID     string  `json:"id"`
	Grades []Grade `json:"grades"` // 新增的字段，存储多个成绩记录
}

var StudentMap map[string]*Student
