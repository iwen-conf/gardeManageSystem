/*
* Editor:KaiWen
* PATH:services/excel_service.go
* Description:
 */

// services/excel_service.go

package services

import (
	"encoding/json"
	"fmt"
	"gardeManageSystem/models"
	"os"
	"strconv"
)

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func AnalyzeGrades(students []*models.Student) map[string]interface{} {
	if len(students) == 0 {
		return nil
	}

	// 计算各科目平均分
	subjectAverages := calculateAverages(students)

	// 计算年级排名
	classRank := calculateRanks(students)

	// 计算进步情况
	improvement := calculateImprovement(students)

	// 计算分数段分布
	distribution := calculateDistribution(students)

	return map[string]interface{}{
		"subjectAverages":     subjectAverages,
		"classRank":           classRank,
		"improvement":         improvement,
		"subjectDistribution": distribution,
	}
}

func calculateAverages(students []*models.Student) map[string]map[string]float64 {
	sums := make(map[string]map[string]float64)
	counts := make(map[string]map[string]int)

	for _, student := range students {
		for _, grade := range student.Grades {
			if sums[grade.ExamName] == nil {
				sums[grade.ExamName] = make(map[string]float64)
				counts[grade.ExamName] = make(map[string]int)
			}

			sums[grade.ExamName]["chinese"] += grade.Chinese
			sums[grade.ExamName]["math"] += grade.Math
			sums[grade.ExamName]["english"] += grade.English
			sums[grade.ExamName]["chemistry"] += grade.Chemistry
			sums[grade.ExamName]["physics"] += grade.Physics
			sums[grade.ExamName]["geography"] += grade.Geography
			sums[grade.ExamName]["biology"] += grade.Biology

			counts[grade.ExamName]["chinese"]++
			counts[grade.ExamName]["math"]++
			counts[grade.ExamName]["english"]++
			counts[grade.ExamName]["chemistry"]++
			counts[grade.ExamName]["physics"]++
			counts[grade.ExamName]["geography"]++
			counts[grade.ExamName]["biology"]++
		}
	}

	averages := make(map[string]map[string]float64)
	for examName, examSums := range sums {
		averages[examName] = make(map[string]float64)
		for subject, sum := range examSums {
			if counts[examName][subject] > 0 {
				averages[examName][subject] = sum / float64(counts[examName][subject])
			}
		}
	}

	return averages
}

func calculateRanks(students []*models.Student) map[string]map[string]int {
	ranks := make(map[string]map[string]int)

	for _, student := range students {
		for _, grade := range student.Grades {
			if ranks[grade.ExamName] == nil {
				ranks[grade.ExamName] = make(map[string]int)
			}
			ranks[grade.ExamName][student.ID]++
		}
	}

	return ranks
}

func calculateImprovement(students []*models.Student) map[string]float64 {
	improvement := make(map[string]float64)

	studentScores := make(map[string][]float64)
	for _, student := range students {
		for _, grade := range student.Grades {
			studentScores[student.ID] = append(studentScores[student.ID], grade.Total)
		}
	}

	for id, scores := range studentScores {
		if len(scores) >= 2 {
			improvement[id] = scores[len(scores)-1] - scores[0]
		}
	}

	return improvement
}

func calculateDistribution(students []*models.Student) map[string]map[string]map[string]int {
	distribution := make(map[string]map[string]map[string]int)
	subjects := []string{"chinese", "math", "english", "chemistry", "physics", "geography", "biology"}

	for _, student := range students {
		for _, grade := range student.Grades {
			if distribution[grade.ExamName] == nil {
				distribution[grade.ExamName] = make(map[string]map[string]int)
			}

			for _, subject := range subjects {
				if distribution[grade.ExamName][subject] == nil {
					distribution[grade.ExamName][subject] = make(map[string]int)
				}

				var score float64
				switch subject {
				case "chinese":
					score = grade.Chinese
				case "math":
					score = grade.Math
				case "english":
					score = grade.English
				case "chemistry":
					score = grade.Chemistry
				case "physics":
					score = grade.Physics
				case "geography":
					score = grade.Geography
				case "biology":
					score = grade.Biology
				}

				switch {
				case score >= 90:
					distribution[grade.ExamName][subject]["90-100"]++
				case score >= 80:
					distribution[grade.ExamName][subject]["80-89"]++
				case score >= 70:
					distribution[grade.ExamName][subject]["70-79"]++
				case score >= 60:
					distribution[grade.ExamName][subject]["60-69"]++
				default:
					distribution[grade.ExamName][subject]["0-59"]++
				}
			}
		}
	}

	return distribution
}

func SaveStudentsToFile() {
	file, err := os.Create("./uploads/students.json")
	if err != nil {
		fmt.Println("Failed to create students.json file:", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(models.StudentMap); err != nil {
		fmt.Println("Failed to encode student data to JSON:", err)
	}
}
func LoadStudentsFromFile() {
	file, err := os.Open("./uploads/students.json")
	if err != nil {
		fmt.Println("Failed to open students.json file:", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&models.StudentMap); err != nil {
		fmt.Println("Failed to decode student data from JSON:", err)
	}
}

func AnalyzeStudentGrades(student *models.Student) map[string]interface{} {
	if student == nil || len(student.Grades) == 0 {
		return nil
	}

	gradesTrend := calculateGradesTrend(student)
	subjectAverages := calculateSubjectAverages(student)
	improvements := calculateImprovements(student)
	strongestAndWeakest := findStrongestAndWeakestSubjects(subjectAverages)

	var examNames []string
	for _, grade := range student.Grades {
		examNames = append(examNames, grade.ExamName)
	}

	return map[string]interface{}{
		"name":             student.Name,
		"id":               student.ID,
		"subjectAverages":  subjectAverages,
		"gradesTrend":      gradesTrend,
		"improvement":      improvements,
		"strongestSubject": strongestAndWeakest["strongest"],
		"weakestSubject":   strongestAndWeakest["weakest"],
		"examNames":        examNames,
	}
}

func calculateSubjectAverages(student *models.Student) map[string]float64 {
	subjectSums := make(map[string]float64)
	subjectCounts := make(map[string]int)

	for _, grade := range student.Grades {
		subjectSums["chinese"] += grade.Chinese
		subjectSums["math"] += grade.Math
		subjectSums["english"] += grade.English
		subjectSums["chemistry"] += grade.Chemistry
		subjectSums["physics"] += grade.Physics
		subjectSums["geography"] += grade.Geography
		subjectSums["biology"] += grade.Biology

		subjectCounts["chinese"]++
		subjectCounts["math"]++
		subjectCounts["english"]++
		subjectCounts["chemistry"]++
		subjectCounts["physics"]++
		subjectCounts["geography"]++
		subjectCounts["biology"]++
	}

	subjectAverages := make(map[string]float64)
	for subject, sum := range subjectSums {
		if count, ok := subjectCounts[subject]; ok && count > 0 {
			subjectAverages[subject] = sum / float64(count)
		}
	}

	return subjectAverages
}

func calculateImprovements(student *models.Student) map[string]float64 {
	improvements := make(map[string]float64)
	if len(student.Grades) < 2 {
		return improvements
	}

	firstGrade := student.Grades[0]
	lastGrade := student.Grades[len(student.Grades)-1]

	improvements["chinese"] = lastGrade.Chinese - firstGrade.Chinese
	improvements["math"] = lastGrade.Math - firstGrade.Math
	improvements["english"] = lastGrade.English - firstGrade.English
	improvements["chemistry"] = lastGrade.Chemistry - firstGrade.Chemistry
	improvements["physics"] = lastGrade.Physics - firstGrade.Physics
	improvements["geography"] = lastGrade.Geography - firstGrade.Geography
	improvements["biology"] = lastGrade.Biology - firstGrade.Biology

	return improvements
}

func findStrongestAndWeakestSubjects(averages map[string]float64) map[string]string {
	var strongest, weakest string
	var gradeMax, gradeMin float64
	for subject, average := range averages {
		if average > gradeMax || strongest == "" {
			gradeMax = average
			strongest = subject
		}
		if average < gradeMin || weakest == "" {
			gradeMin = average
			weakest = subject
		}
	}
	return map[string]string{
		"strongest": strongest,
		"weakest":   weakest,
	}
}

func calculateGradesTrend(student *models.Student) map[string][]float64 {
	gradesTrend := make(map[string][]float64)
	for _, grade := range student.Grades {
		gradesTrend["chinese"] = append(gradesTrend["chinese"], grade.Chinese)
		gradesTrend["math"] = append(gradesTrend["math"], grade.Math)
		gradesTrend["english"] = append(gradesTrend["english"], grade.English)
		gradesTrend["chemistry"] = append(gradesTrend["chemistry"], grade.Chemistry)
		gradesTrend["physics"] = append(gradesTrend["physics"], grade.Physics)
		gradesTrend["geography"] = append(gradesTrend["geography"], grade.Geography)
		gradesTrend["biology"] = append(gradesTrend["biology"], grade.Biology)
		gradesTrend["total"] = append(gradesTrend["total"], grade.Total) // 确保总分存在
	}
	return gradesTrend
}
