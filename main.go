package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Task represents a task in the database
type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

func main() {
	db, err := gorm.Open("sqlite3", "tasks.db")
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	// Auto-migrate the schema
	db.AutoMigrate(&Task{})

	r := gin.Default()
	
	// just to check
	r.GET("/check", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// this route will create a new task
	r.POST("/tasks", func(c *gin.Context) {
		var task Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&task)
		c.JSON(http.StatusCreated, task)
	})

	// retrieve task by ID
	r.GET("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusOK, task)
	})

	// update a task by ID
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&task)
		c.JSON(http.StatusOK, task)
	})

	

	// delete a task by ID
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		db.Delete(&task)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	})

	// list all tasks
	r.GET("/tasks", func(c *gin.Context) {
		var tasks []Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})

	// run the server
	port := ":8080"
	log.Printf("Server is listening on port %s", port)
	r.Run(port)
}
