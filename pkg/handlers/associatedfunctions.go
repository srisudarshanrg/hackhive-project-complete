package handlers

import (
	"log"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(hashedPassword), nil
}

// GetPasswordFromHash compares a password to its hash
func GetPasswordFromHash(entered_password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(entered_password))
	return err == nil
}

// SendEmail sends an email to the given email
func SendEmail(from string, to []string, message []byte, password string) error {
	// smtp server configurations
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// SetOtp sets the otp value so that is accessible throughout the file
func SetOtp(otp_recieved string) {
	otp = otp_recieved
}

// SetLoginStatus sets the login status so that every handler can access it
func SetLoginStatus(loginStatusArg bool) {
	loginStatus = loginStatusArg
}
