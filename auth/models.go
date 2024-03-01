package auth

type User struct {
    ID              int    `json:"id"`
    Username        string `json:"username"`
    Password        string `json:"-"`
    Email           string `json:"email"`
    TOTPSecret      string `json:"-"`
    IsEmailVerified bool   `json:"isEmailVerified"`
    MFAEnabled      bool   `json:"mfaEnabled"`
}
