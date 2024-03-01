package auth

import (
	"educational_api/db"
	"educational_api/models"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var smtpHost = os.Getenv("SMTP_HOST")
var smtpPort = os.Getenv("SMTP_PORT")
var smtpUser = os.Getenv("SMTP_USER")
var smtpPass = os.Getenv("SMTP_PASS")


func sendVerificationEmail(to, token string) error {

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration not set")
	}

	from := smtpUser
	subject := "Email Verification"
	body := fmt.Sprintf("Here is your verification token: %s", token)
	msg := createMessage(from, to, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
			return err
	}
	return nil
}

func createMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return []byte(message)
}
func saveUser(user models.User) error {
	fmt.Println("Saving user")

	stmt, err := db.DB.Prepare("INSERT INTO users(username, password, email, totpSecret, isEmailVerified, mfaEnabled) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
			return fmt.Errorf("error preparing user insert: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password, user.Email, user.TOTPSecret, user.IsEmailVerified, user.MFAEnabled)
	if err != nil {
			return fmt.Errorf("error executing user insert: %w", err)
	}

	return nil
}

func emailExists(email string) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
			log.Println("Error checking if email exists:", err)
	}
	return exists
}

func updateUserEmailVerificationStatus(email string, isVerified bool) error {
	stmt, err := db.DB.Prepare("UPDATE users SET isEmailVerified = ? WHERE email = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(isVerified, email)
	if err != nil {
		return err
	}

	return nil
}
