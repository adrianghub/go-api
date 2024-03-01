package auth

import (
	"educational_api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	fmt.Println("Signup handler called")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	user.IsEmailVerified = false
	user.MFAEnabled = false

	if emailExists(user.Email) {
    http.Error(w, "Email already in use", http.StatusBadRequest)
    return
	}

	if err := saveUser(user); err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	verificationToken, err := GenerateEmailVerificationToken(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	err = sendVerificationEmail(user.Email, verificationToken)
	if err != nil {
			log.Fatalf("Failed to send verification email: %v", err)
			return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully. Please check your email to verify."})
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
	}

	if err := VerifyUserToken(token); err != nil {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
	}

	// Redirect or respond with success
	// http.Redirect(w, r, "/verification-success", http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
}