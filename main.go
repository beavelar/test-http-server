package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var logger *log.Logger

type User struct {
	ID       string
	Username string
}

type Message struct {
	ID      string
	UserID  string
	Message string
}

func init() {
	logger = log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:   time.DateTime,
		Level:        log.DebugLevel,
		ReportCaller: true,
		Formatter:    log.JSONFormatter,
	})
}

func main() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		logger.Fatal("DATABASE_URL environment variable not set. Please set it to your PostgreSQL connection string.")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Fatal("Error opening database", "error", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal("Error connecting to the database", "error", err)
	}
	logger.Info("Successfully connected to the database!")

	router := gin.Default()
	router.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, username FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			user := User{}
			if err := rows.Scan(&user.ID, &user.Username); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		res := []gin.H{}
		for _, user := range users {
			res = append(res, gin.H{
				"id":       user.ID,
				"username": user.Username,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"users": res,
		})
	})

	router.GET("/users/:id/messages", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("SELECT message_id, user_id, message FROM messages WHERE user_id = ?")
		if err != nil {
			logger.Fatal("Failed to prepare statement for users", "error", err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		messages := []Message{}
		for rows.Next() {
			msg := Message{}
			if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Message); err != nil {
				log.Fatal(err)
			}
			messages = append(messages, msg)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		res := []gin.H{}
		for _, msg := range messages {
			res = append(res, gin.H{
				"id":       msg.ID,
				"userId":   msg.UserID,
				"messaage": msg.Message,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": res,
		})
	})

	router.Run()
}
