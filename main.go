package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWT struct {
	ID      int    `json:"id"`
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func postUsers(c *gin.Context) {
	var newUser User

	// Calls BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// Adds the new album to the slice.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

var users = []User{
	{ID: 1, Email: "test@example.com", Password: "sdf211sd"},
}

func main() {
	fmt.Println("users", users)
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", postUsers)

	router.Run("localhost:8080")
}
