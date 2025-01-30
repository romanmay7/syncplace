package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/romanmay7/syncplace/wsocket"
	//"github.com/google/uuid"
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
	GetRoomChatMessages(roomId string) ([]wsocket.ChatMessage, error)
	AddNewChatMessage(roomId string, msgID string, timestamp string, content string, sender string, file_path string) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable not set")
	}

	//connStr := "user=postgres dbname=postgres password=syncplace sslmode=disable"

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

	fmt.Println("calling : CreateAppTables()")

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

	if err != nil {
		return err
	}

	query3 := `CREATE TABLE IF NOT EXISTS chat_message (
		id SERIAL PRIMARY KEY,
		msg_id  UUID NOT NULL UNIQUE,
		room_id UUID NOT NULL,
        time_stamp timestamp,
		content TEXT,
		sender varchar(50),
		file_path TEXT
    );`

	_, err = s.db.Exec(query3)
	fmt.Println(err)

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

	fmt.Println("calling: GetBoardState for Room: " + roomId)

	var elements []interface{}

	rows, err := s.db.Query("select data from board_state where room_id = $1", roomId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		fmt.Println("Deserializing data")
		var jsonData []byte
		err := rows.Scan(&jsonData)
		if err != nil {
			return nil, err
		}

		// Deserializing elements data into generic slice (no specific structure)
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

// ------------------------------------------------------------------------------------------------------------------

func (s *PostgresStore) GetRoomChatMessages(roomId string) ([]wsocket.ChatMessage, error) {
	fmt.Println("calling: GetRoomChatMessages")
	var chatMessages []wsocket.ChatMessage

	rows, err := s.db.Query("select msg_id, room_id, time_stamp, content, sender, file_path from chat_message where room_id = $1", roomId)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Close the rows after iterating

	for rows.Next() {
		var msgId string
		var roomId string
		var timestamp time.Time
		var content string
		var sender string
		var file_path string

		err := rows.Scan(&msgId, &roomId, &timestamp, &content, &sender, &file_path)
		if err != nil {
			return nil, err
		}

		// Create the chat message object with the scanned values
		chatMessage := wsocket.ChatMessage{
			MsgID:     msgId,
			RoomID:    roomId,
			Timestamp: timestamp.String(),
			Content:   content,
			Sender:    sender,
			FilePath:  file_path,
		}

		chatMessages = append(chatMessages, chatMessage)

	}

	return chatMessages, err
}

func (s *PostgresStore) AddNewChatMessage(roomId string, msgID string, timestamp string, content string, sender string, file_path string) error {

	query := `insert into chat_message
	        (msg_id, room_id, time_stamp, content, sender, file_path)
			 values ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (msg_id) DO NOTHING`

	fmt.Println("INSERT:" + msgID + "|" + roomId + "|" + timestamp + "|" + content + "|" + sender + "|" + file_path)
	//msgID := uuid.New()
	res, err := s.db.Query(
		query,
		msgID,
		roomId,
		timestamp,
		content,
		sender,
		file_path)

	if err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil

}

// ------------------------------------------------------------------------------------------------------------------
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
