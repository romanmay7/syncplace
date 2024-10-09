package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

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

func (a *UserAccount) ValidatePassword(psw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(psw)) == nil
}

func NewUserAccount(username string, email string, password string) (*UserAccount, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &UserAccount{
		//ID:        rand.Intn(10000),
		UserName:  username,
		Email:     email,
		Password:  string(encpw), //Encrypted Password
		CreatedAt: time.Now().UTC(),
	}, nil

}
