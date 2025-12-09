package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"regexp"
)

// SendEmail sends an email using Gmail SMTP server
func SendEmail(toEmail, subject, body string) error {
	
	from := os.Getenv("EMAIL_FROM")         // Gmail email from environment variable
	password := os.Getenv("EMAIL_PASSWORD") // Gmail app password from env variable
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := []byte("Subject: " + subject + "\n\n" + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
	if err != nil {
		return err
	}
	fmt.Println("✅ Email sent to:", toEmail)
	return nil
}
// ValidateEmail checks basic email format
func ValidateEmail(email string) bool {
	// simple regex check
	re := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched := regexp.MustCompile(re).MatchString(email)
	return matched
}
