package main

import "math/rand"

type UserAccount struct {
	ID       int    `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserAccount(username string, email string, password string) *UserAccount {
	return &UserAccount{
		ID:       rand.Intn(10000),
		UserName: username,
		Email:    email,
		Password: password,
	}

}
