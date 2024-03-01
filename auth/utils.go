package auth

import (
	"bufio"
	"database/sql"
	"educational_api/db"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/mdp/qrterminal"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

func sendVerificationEmail(to, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration not set")
	}

	from := smtpUser
	subject := "Email Verification"
	body := fmt.Sprintf("Click <a href='http://localhost:8080/verify-email?token=%s'>here</a> to verify your email", token)

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

func saveUser(user User) error {
	_, err := db.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, email TEXT, totpSecret TEXT, isEmailVerified BOOLEAN, mfaEnabled BOOLEAN)")
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

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

func getUserByEmail(userEmail string) (*User, error) {
	var user User
	query := `SELECT id, username, password, email, totpSecret, isEmailVerified, mfaEnabled FROM users WHERE email = $1`
	
	err := db.DB.QueryRow(query, userEmail).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.TOTPSecret, &user.IsEmailVerified, &user.MFAEnabled)
	if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("user not found or email not verified") 
			}
			return nil, err
	}

	if !user.IsEmailVerified {
		return nil, errors.New("user not found or email not verified")
	}

	return &user, nil
}

func generateTOTPWithSecret(randomSecret, userEmail string) {
	uri := gotp.NewDefaultTOTP(randomSecret).ProvisioningUri(userEmail, "myApp")
	fmt.Println("Secret Key URI:", uri)
 
	qrcode.WriteFile(uri, qrcode.Medium, 256, "qr.png")
 
	qrterminal.GenerateWithConfig(uri, qrterminal.Config{
	 Level:     qrterminal.L,
	 Writer:    os.Stdout,
	 BlackChar: qrterminal.BLACK,
	 WhiteChar: qrterminal.WHITE,
	})
 
	fmt.Println("\nScan the QR code with your authenticator app")
 }

 func verifyOTP(randomSecret string) string {
	totp := gotp.NewDefaultTOTP(randomSecret)

 fmt.Print("Enter the OTP from your authenticator app: ")
 
 scanner := bufio.NewScanner(os.Stdin)
 
 scanner.Scan()
 
 userInput := scanner.Text()

	if totp.Verify(userInput, time.Now().Unix()) {
		fmt.Println("Authentication successful! Access granted.")
	} else {
		fmt.Println("Authentication failed! Invalid OTP.")
		return ""
	}

	return userInput
}