package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	fmt.Println("Signup handler called")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(user.Password)

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

	verificationToken, err := generateEmailVerificationToken(user.Email)
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

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
	}

	if err := verifyEmailVerificationToken(token); err != nil {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
	}

	// Redirect or respond with success
	// http.Redirect(w, r, "/verification-success", http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
	}

	user, err := getUserByEmail(credentials.Email)

	if err != nil || !checkPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(map[string]string{"message": "Authentication successful"})
}

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/signup", signUpHandler).Methods("POST")
	router.HandleFunc("/verify-email", verifyEmailHandler).Methods("GET")
	router.HandleFunc("/signin", signInHandler).Methods("POST")
}