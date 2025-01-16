package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexbrainman/printer"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// Models
type Test struct {
	ID        float64    `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID            int      `json:"id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
}

type TestResult struct {
	UserID     string      `json:"userId"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	MiddleName string      `json:"middleName"`
	TestID     float64     `json:"testId"`
	Answers    map[int]int `json:"answers"`
	Score      int         `json:"score"`
	Date       time.Time   `json:"date"`
}

type TestResultResponse struct {
	UserID     string    `json:"userId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	MiddleName string    `json:"middleName"`
	TestID     float64   `json:"testId"`
	Score      int       `json:"score"`
	Percentage float64   `json:"percentage"`
	Date       time.Time `json:"date"`
}

type PrintTask struct {
	TaskID    string
	UserID    string
	TestID    float64
	Status    string
	Timestamp time.Time
	Error     string
}

var (
	printQueue   = make(chan PrintTask, 100)
	taskStatuses = sync.Map{}
	taskLock     sync.Mutex
)

var tests []Test

func main() {
	// Load test data
	loadTestData()

	go processPrintQueue()

	// Initialize Gin router
	r := gin.Default()

	// Use CORS middleware
	r.Use(corsMiddleware())

	// Define routes
	r.GET("/api/tests", getTests)
	r.GET("/api/tests/:id", getTestById)
	r.POST("/api/results", submitTestResult)
	r.GET("/api/results/:userId/:testId", getTestResultById)
	r.GET("/api/generate-pdf/:userId/:testId", func(c *gin.Context) {
		userId := c.Param("userId")
		testId := c.Param("testId")

		testIdFloat, err := strconv.ParseFloat(testId, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid test ID format"})
			return
		}

		pdfPath, err := generatePDF(userId, testIdFloat)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to generate PDF: %v", err)})
			return
		}

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=test_result_%s_%v.pdf", userId, testId))
		c.Header("Content-Type", "application/pdf")
		c.File(pdfPath)

		// Удалим временный файл после отправки
		go func() {
			time.Sleep(5 * time.Second) // Даем время на передачу файла
			os.Remove(pdfPath)
		}()
	})

	// Start server
	r.Run(":8080")
}

func getTestResultById(c *gin.Context) {
	userID := c.Param("userId")
	testID := c.Param("testId")

	testIDFloat, err := strconv.ParseFloat(testID, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid test ID format"})
		return
	}

	files, err := os.ReadDir("temp")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read results directory"})
		return
	}

	var result TestResult
	var found bool
	var latestTime int64

	// Find the most recent result for the given user and test
	for _, file := range files {
		if strings.Contains(file.Name(), fmt.Sprintf("%s_%v", userID, testIDFloat)) {
			data, err := os.ReadFile(filepath.Join("temp", file.Name()))
			if err != nil {
				continue
			}

			var tempResult TestResult
			if err := json.Unmarshal(data, &tempResult); err != nil {
				continue
			}

			if tempResult.Date.Unix() > latestTime {
				latestTime = tempResult.Date.Unix()
				result = tempResult
				found = true
			}
		}
	}

	if !found {
		c.JSON(404, gin.H{"error": "Test result not found"})
		return
	}

	score, percentage := calculateScore(result)

	response := TestResultResponse{
		UserID:     result.UserID,
		FirstName:  result.FirstName,
		LastName:   result.LastName,
		MiddleName: result.MiddleName,
		TestID:     result.TestID,
		Score:      score,
		Percentage: percentage,
		Date:       result.Date,
	}

	c.JSON(200, response)
}

func loadTestData() {
	data, err := os.ReadFile("test-questions.json")
	if err != nil {
		panic(err)
	}

	var testData struct {
		Tests []Test `json:"tests"`
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		panic(err)
	}

	tests = testData.Tests
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// API Handlers
func getTests(c *gin.Context) {
	var testsWithoutAnswers []Test

	for _, test := range tests {
		newTest := Test{
			ID:    test.ID,
			Title: test.Title,
		}

		for _, question := range test.Questions {
			newQuestion := Question{
				ID:       question.ID,
				Question: question.Question,
				Options:  question.Options,
			}
			newTest.Questions = append(newTest.Questions, newQuestion)
		}

		testsWithoutAnswers = append(testsWithoutAnswers, newTest)
	}

	c.JSON(200, gin.H{
		"tests": testsWithoutAnswers,
	})
}

func getTestById(c *gin.Context) {
	testId := c.Param("id")

	for _, test := range tests {
		if fmt.Sprint(test.ID) == testId {
			// Create test without correct answers
			sanitizedTest := Test{
				ID:    test.ID,
				Title: test.Title,
			}

			for _, question := range test.Questions {
				sanitizedQuestion := Question{
					ID:       question.ID,
					Question: question.Question,
					Options:  question.Options,
				}
				sanitizedTest.Questions = append(sanitizedTest.Questions, sanitizedQuestion)
			}

			c.JSON(200, sanitizedTest)
			return
		}
	}

	c.JSON(404, gin.H{"error": "Test not found"})
}

func processPrintQueue() {
	taskLock.Lock()
	defer taskLock.Unlock()

	for task := range printQueue {
		if task.Status == "completed" || task.Status == "failed" {
			log.Printf("Задача %v уже завершена, пропускаем.", task.TaskID)
			continue
		}

		log.Printf("Начало обработки задачи: %v", task.TaskID)
		task.Status = "processing"
		task.Timestamp = time.Now()
		taskStatuses.Store(task.TaskID, task)

		pdfPath, err := generatePDF(task.UserID, task.TestID)
		if err != nil {
			log.Printf("Ошибка генерации PDF для задачи %s: %v", task.TaskID, err)
			task.Status = "failed"
			task.Error = err.Error()
			taskStatuses.Store(task.TaskID, task)
			continue
		}

		err = sendToPrinterWithRetries(pdfPath)
		if err != nil {
			log.Printf("Ошибка печати для задачи %s: %v", task.TaskID, err)
			task.Status = "failed"
			task.Error = err.Error()
			taskStatuses.Store(task.TaskID, task)
			continue
		}

		log.Printf("Успешная обработка задачи: %v", task.TaskID)
		task.Status = "completed"
		task.Timestamp = time.Now()
		taskStatuses.Store(task.TaskID, task)
	}
}

func submitTestResult(c *gin.Context) {
	var result TestResult
	if err := c.BindJSON(&result); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	result.Date = time.Now()

	score, percentage := calculateScore(result)
	result.Score = score

	filename := fmt.Sprintf("temp/%s_%v_%d.json", result.UserID, result.TestID, result.Date.Unix())

	os.MkdirAll("temp", 0755)
	file, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process result"})
		return
	}

	if err := os.WriteFile(filename, file, 0644); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save result"})
		return
	}

	taskID := fmt.Sprintf("%s_%v_%d", result.UserID, result.TestID, time.Now().Unix())
	task := PrintTask{
		TaskID:    taskID,
		UserID:    result.UserID,
		TestID:    result.TestID,
		Status:    "created",
		Timestamp: time.Now(),
	}

	printQueue <- task

	c.JSON(200, gin.H{
		"userId":     result.UserID,
		"firstName":  result.FirstName,
		"lastName":   result.LastName,
		"testId":     result.TestID,
		"score":      result.Score,
		"percentage": fmt.Sprintf("%.2f%%", percentage),
		"date":       result.Date,
	})
}

func calculateScore(result TestResult) (int, float64) {
	var correctAnswers int
	var test Test

	// Find the test
	for _, t := range tests {
		if t.ID == result.TestID {
			test = t
			break
		}
	}

	for _, q := range test.Questions {
		if answer, exists := result.Answers[q.ID]; exists {
			if answer == q.CorrectAnswer {
				correctAnswers++
			}
		}
	}

	totalQuestions := len(test.Questions)
	percentage := 0.0

	if totalQuestions > 0 {
		percentage = (float64(correctAnswers) / float64(totalQuestions)) * 100
	}

	return correctAnswers, percentage
}

type TestScore struct {
	Score      int     // Количество правильных ответов
	Percentage float64 // Процент правильных ответов
}

func generatePDF(userID string, testID float64) (string, error) {
	files, err := os.ReadDir("temp")
	if err != nil {
		return "", fmt.Errorf("failed to read results directory: %w", err)
	}

	var result TestResult
	var found bool

	var latestTime int64
	for _, file := range files {
		if strings.Contains(file.Name(), fmt.Sprintf("%s_%v", userID, testID)) {
			data, err := os.ReadFile(filepath.Join("temp", file.Name()))
			if err != nil {
				continue
			}

			var tempResult TestResult
			if err := json.Unmarshal(data, &tempResult); err != nil {
				continue
			}

			if tempResult.Date.Unix() > latestTime {
				latestTime = tempResult.Date.Unix()
				result = tempResult
				found = true
			}
		}
	}

	if !found {
		return "", fmt.Errorf("test result not found for user %s and test %v", userID, testID)
	}

	// Находим информацию о тесте
	var test Test
	for _, t := range tests {
		if t.ID == testID {
			test = t
			break
		}
	}

	// Рассчитываем результаты
	score, percentage := calculateScore(result)

	pdfPath := fmt.Sprintf("results/%s_%v.pdf", userID, testID)
	os.MkdirAll("results", 0755)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	fontPath := "./asset/font/Roboto-VariableFont_wdth,wght.ttf"
	pdf.AddUTF8Font("Roboto", "", fontPath)
	pdf.AddUTF8Font("Roboto", "B", fontPath)

	// Заголовок
	pdf.SetFont("Roboto", "B", 20)
	pdf.SetTextColor(0, 51, 102)
	pdf.CellFormat(190, 20, "Результаты тестирования", "0", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Информация о тесте
	pdf.SetFont("Roboto", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 8, fmt.Sprintf("Название теста: %s", test.Title), "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	// Информация о пользователе
	pdf.SetFont("Roboto", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("ФИО: %s %s %s", result.LastName, result.FirstName, result.MiddleName), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Дата прохождения: %s", result.Date.Format("02.01.2006 15:04")), "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	// Результаты
	pdf.SetFont("Roboto", "B", 14)
	pdf.SetTextColor(0, 102, 51)
	pdf.CellFormat(190, 10, "Результаты:", "0", 1, "L", false, 0, "")

	pdf.SetFont("Roboto", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 8, fmt.Sprintf("Правильных ответов: %d из %d", score, len(test.Questions)), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Процент выполнения: %.1f%%", percentage), "0", 1, "L", false, 0, "")

	// Оценка результата
	pdf.Ln(10)
	pdf.SetFont("Roboto", "B", 14)
	var resultText string
	switch {
	case percentage >= 90:
		resultText = "Отлично"
		pdf.SetTextColor(0, 128, 0)
	case percentage >= 80:
		resultText = "Хорошо"
		pdf.SetTextColor(0, 128, 128)
	default:
		resultText = "Требуется повторное прохождение"
		pdf.SetTextColor(255, 0, 0)
	}
	pdf.CellFormat(190, 10, fmt.Sprintf("Итоговая оценка: %s", resultText), "0", 1, "L", false, 0, "")

	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	taskID := fmt.Sprintf("%s_%v_%d", userID, testID, time.Now().Unix())
	task := PrintTask{
		TaskID:    taskID,
		UserID:    userID,
		TestID:    testID,
		Status:    "created",
		Timestamp: time.Now(),
	}
	taskStatuses.Store(taskID, task)
	printQueue <- task

	return pdfPath, nil
}

// Максимальное количество попыток отправки на печать
const maxAttempts = 5

// Интервал между попытками (в секундах)
const retryInterval = 2

func sendToPrinterWithRetries(pdfPath string) error {
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err := sendToPrinter(pdfPath)
		if err == nil {
			return nil
		}

		log.Printf("Попытка %d из %d завершилась ошибкой: %v", attempts, maxAttempts, err)
		time.Sleep(time.Duration(attempts*retryInterval) * time.Second)
	}
	return fmt.Errorf("все попытки отправки на печать завершились неудачей")
}

func sendToPrinter(pdfPath string) error {
	printerName, err := printer.Default()
	if err != nil {
		return fmt.Errorf("не удалось получить имя принтера: %w", err)
	}

	log.Printf("Отправка PDF на принтер %s: %s", printerName, pdfPath)

	fileData, err := os.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения PDF: %w", err)
	}

	p, err := printer.Open(printerName)
	if err != nil {
		return fmt.Errorf("не удалось открыть принтер: %w", err)
	}
	defer p.Close()

	if err := p.StartDocument("Test Print", "RAW"); err != nil {
		return fmt.Errorf("не удалось начать документ: %w", err)
	}
	defer p.EndDocument()

	if _, err := p.Write(fileData); err != nil {
		return fmt.Errorf("ошибка записи данных в принтер: %w", err)
	}

	log.Printf("Файл успешно отправлен на печать: %s", pdfPath)
	return nil
}
