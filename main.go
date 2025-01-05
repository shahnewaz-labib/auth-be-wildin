package main

import (
	"auth-be-wildin/auth"
	"auth-be-wildin/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	PORT = ":9000"
)

func init() {
	db.InitDB()
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerHandler(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid format"})
		return
	}

	err := auth.RegisterUser(loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func loginHandler(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid format"})
		return
	}

	isAuthenticated, err := auth.AuthenticateUser(loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	sessionId := uuid.New().String()

	err = auth.CreateSession(sessionId, loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create session"})
		return
	}

	// key, value, maxAge, path, domain, secure, httpOnly
	c.SetCookie("sessionId", sessionId, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func authMiddleware(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		c.Abort()
		return
	}

	_, err = auth.GetSessionUsername(sessionId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid session"})
		c.Abort()
		return
	}

	c.Next()
}

func getLoggedInUser(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no user logged in"})
		return
	}

	username, err := auth.GetSessionUsername(sessionId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username})
}

func logoutHandler(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err == nil {
		err = auth.DeleteSession(sessionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete session"})
			return
		}
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

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	auth := r.Group("/")
	auth.Use(authMiddleware)
	{
		auth.GET("/me", getLoggedInUser)
		auth.POST("/logout", logoutHandler)
	}

	r.Run(PORT)
}
