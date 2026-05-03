package models

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
}

type UserWithPassword struct {
	Email        string
	Name         string
	PasswordHash string
	Id           int
}
