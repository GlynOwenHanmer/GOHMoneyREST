package model

type User struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
