package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Struct สำหรับ map ข้อมูล
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	var err error
	// เชื่อมต่อ MySQL
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	// เปิดใช้งาน CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/users", getUsers)

	r.GET("/users/:id", getUserByID)

	r.Run(":8080")
}

// ฟังก์ชันดึงข้อมูลทั้งหมด
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// ฟังก์ชันดึงข้อมูลตาม ID
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	var user User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
