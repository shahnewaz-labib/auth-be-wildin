package auth

import (
	"auth-be-wildin/db"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.GetDB().Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func AuthenticateUser(username, password string) (bool, error) {
	var hashedPassword string
	err := db.GetDB().QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func CreateSession(sessionId, username string) error {
	_, err := db.GetDB().Exec("INSERT INTO sessions (sessionId, username) VALUES (?, ?)", sessionId, username)

	return err
}

func GetSessionUsername(sessionId string) (string, error) {
	var username string
	err := db.GetDB().QueryRow("SELECT username FROM sessions WHERE sessionId = ?", sessionId).Scan(&username)

	return username, err
}

func DeleteSession(sessionId string) error {
	_, err := db.GetDB().Exec("DELETE FROM sessions WHERE sessionId = ?", sessionId)

	return err
}
