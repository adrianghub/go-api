package auth

type User struct {
    ID              int    `json:"id"`
    Username        string `json:"username"`
    Password        string `json:"password"`
    Email           string `json:"email"`
    TOTPSecret      string `json:"-"`
    IsEmailVerified bool   `json:"isEmailVerified"`
    MFAEnabled      bool   `json:"mfaEnabled"`
}

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}