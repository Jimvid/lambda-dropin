package model

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
