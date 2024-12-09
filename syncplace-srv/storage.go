package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUserAccount(*UserAccount) error
	DeleteUserAccount(int) error
	UpdateUserAccount(*UserAccount) error
	GetUserAccounts() ([]*UserAccount, error)
	GetUserAccountByID(int) (*UserAccount, error)
	GetUserAccountByUserName(string) (*UserAccount, error)
	GetBoardState(string) ([]interface{}, error)
	CreateBoardStateRecord(string, string, []interface{}) error
	UpdateBoardStateRecord(string, []interface{}) error
	CheckIfBoardRecordExist(string) (bool, error)
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
	return s.CreateAppTables()
}

func (s *PostgresStore) CreateAppTables() error {
	query := `create table if not exists user_account(
     id serial primary key,
	 user_name varchar(50),
	 email varchar(50),
	 password varchar(100),
	 created_at timestamp
   )`

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	query2 := `CREATE TABLE IF NOT EXISTS board_state (
			id SERIAL PRIMARY KEY,
			room_id UUID NOT NULL,
			room_name varchar(50),
			data JSONB NOT NULL
	);
`
	_, err = s.db.Exec(query2)

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

func (s *PostgresStore) GetBoardState(roomId string) ([]interface{}, error) {

	var elements []interface{}

	rows, err := s.db.Query("select data from board_state where room_id = $1", roomId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var jsonData []byte
		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		// Deserializing elements data into generic slice (if no specific structure)
		var elements []interface{}
		err = json.Unmarshal(jsonData, &elements)

		if err != nil {
			fmt.Print(err)
		}

	}

	//return nil, fmt.Errorf("Room State %d not found", id)

	return elements, err

}

func (s *PostgresStore) CreateBoardStateRecord(roomId string, roomName string, elements []interface{}) error {

	query := `insert into board_state
	        (room_id, room_name, data)
			 values ($1, $2, $3)`

	// Convert elements to JSON
	elementsJsonData, err := json.Marshal(elements)
	if err != nil {
		return err
	}

	res, err := s.db.Query(
		query,
		roomId,
		roomName,
		elementsJsonData)

	if err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil

}

func (s *PostgresStore) CheckIfBoardRecordExist(roomId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM board_state WHERE room_id = $1)`

	rows, err := s.db.Query(query, roomId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var exists bool
	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil // No rows found
	}

	return exists, nil
}

func (s *PostgresStore) UpdateBoardStateRecord(roomId string, elements []interface{}) error {
	query := `UPDATE board_state
			   SET data = $1
			   WHERE room_id = $2`

	_, err := s.db.Exec(query, elements, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateUserAccount(*UserAccount) error {
	return nil
}

func (s *PostgresStore) DeleteUserAccount(id int) error {
	_, err := s.db.Query("delete from user_account where id = $1", id)

	return err
}

func scanRowsIntoUserAccount(rows *sql.Rows) (*UserAccount, error) {
	account := new(UserAccount)
	err := rows.Scan(
		&account.ID,
		&account.UserName,
		&account.Email,
		&account.Password,
		&account.CreatedAt)

	return account, err
}

func (s *PostgresStore) GetUserAccountByID(id int) (*UserAccount, error) {
	rows, err := s.db.Query("select * from user_account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanRowsIntoUserAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetUserAccountByUserName(username string) (*UserAccount, error) {
	rows, err := s.db.Query("select * from user_account where user_name = $1", username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanRowsIntoUserAccount(rows)
	}

	return nil, fmt.Errorf("Account %s not found", username)
}

func (s *PostgresStore) GetUserAccounts() ([]*UserAccount, error) {
	rows, err := s.db.Query("select * from user_account")

	if err != nil {
		return nil, err
	}

	usrAccounts := []*UserAccount{}

	for rows.Next() {
		account, err := scanRowsIntoUserAccount(rows)

		if err != nil {
			return nil, err
		}

		usrAccounts = append(usrAccounts, account)
	}

	return usrAccounts, nil
	//fmt.Printf("%+v\n", rows)

}
