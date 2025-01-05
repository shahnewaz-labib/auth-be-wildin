package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	PORT = ":9000"
)

var sessions = make(map[string]string)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func validateCredentials(username string, password string) bool {
	return true
}

func loginHandler(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid format"})
		return
	}

	if !validateCredentials(loginReq.Username, loginReq.Username) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	sessionId := uuid.New().String()

	sessions[sessionId] = loginReq.Username

	// key, value, maxAge, path, domain, secure, httpOnly
	c.SetCookie("sessionId", sessionId, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	printAllSessions()
}

func printAllSessions() {
	log.Println("Current sessions:")
	for sessionID, username := range sessions {
		log.Printf("SessionID: %s, Username: %s\n", sessionID, username)
	}
}

func authMiddleware(c *gin.Context) {
	printAllSessions()
	sessionId, err := c.Cookie("sessionId")

	if err != nil || sessions[sessionId] == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}

func getLoggedInUser(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
		return
	}

	username := sessions[sessionId]
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username})
}

func logoutHandler(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err == nil {
		delete(sessions, sessionId)
	}

	c.SetCookie("sessionId", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	r := gin.Default()

	r.GET("/ping", pingHandler)

	r.POST("/login", loginHandler)

	auth := r.Group("/")
	auth.Use(authMiddleware)
	{
		auth.GET("/me", getLoggedInUser)
		auth.POST("/logout", logoutHandler)
	}

	r.Run(PORT)
}
