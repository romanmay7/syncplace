package main

import (
	"time"
)

type UserAccount struct {
	ID        int       `json:"id"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateUserAccountRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewUserAccount(username string, email string, password string) *UserAccount {
	return &UserAccount{
		//ID:        rand.Intn(10000),
		UserName:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}

}
