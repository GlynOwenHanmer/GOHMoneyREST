package model

// User contains the details of a user.
type User struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
