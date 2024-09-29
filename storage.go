package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUserAccount(*UserAccount) error
	DeleteUserAccount(int) error
	UpdateUserAccount(*UserAccount) error
	GetUserAccounts() ([]*UserAccount, error)
	GetUserAccountByID(int) (*UserAccount, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=syncplace sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

func (s *PostgresStore) Init() error {
	return s.CreateUserAccountTable()
}

func (s *PostgresStore) CreateUserAccountTable() error {
	query := `create table if not exists user_account(
     id serial primary key,
	 user_name varchar(50),
	 email varchar(50),
	 password varchar(50),
	 created_at timestamp
   )`

	_, err := s.db.Exec(query)
	return err
}

//=============Implement 'Storage' Interface====================================

func (s *PostgresStore) CreateUserAccount(acc *UserAccount) error {
	query := `insert into user_account
	        (user_name, email, password, created_at)
			 values ($1, $2, $3, $4)`

	res, err := s.db.Query(
		query,
		acc.UserName,
		acc.Email,
		acc.Password,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil
}

func (s *PostgresStore) UpdateUserAccount(*UserAccount) error {
	return nil
}

func (s *PostgresStore) DeleteUserAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetUserAccountByID(id int) (*UserAccount, error) {
	return nil, nil
}

func (s *PostgresStore) GetUserAccounts() ([]*UserAccount, error) {
	rows, err := s.db.Query("select * from user_account")

	if err != nil {
		return nil, err
	}

	usrAccounts := []*UserAccount{}

	for rows.Next() {
		account := new(UserAccount)
		err := rows.Scan(
			&account.ID,
			&account.UserName,
			&account.Email,
			&account.Password,
			&account.CreatedAt)

		if err != nil {
			return nil, err
		}

		usrAccounts = append(usrAccounts, account)
	}

	return usrAccounts, nil
	//fmt.Printf("%+v\n", rows)

}
