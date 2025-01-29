package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// Mock Store and other dependencies as needed
type MockStore struct{} // to Implement  Store interface methods

type MockFileManager struct{}

func (m *MockFileManager) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Mock implementation ( always return success)
	w.WriteHeader(http.StatusOK)
}

type MockWsockHandler struct{}

// Defining Test Server struct with all the fields
type Server struct {
	store        MockStore
	filemgr      MockFileManager
	wsockHandler *MockWsockHandler
}

type RoomsMockResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientsMockResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (m *MockWsockHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m *MockWsockHandler) GetRooms(w http.ResponseWriter, r *http.Request) {

	rooms := make([]RoomsMockResult, 0)

	room := &RoomsMockResult{
		ID:   "123321342142412421",
		Name: "Test Room",
	}

	rooms = append(rooms, *room)

	WriteJSON(w, http.StatusOK, rooms)

}

// JoinRoom method--------------------------------------------------------------------------------
func (m *MockWsockHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m *MockWsockHandler) GetClients(w http.ResponseWriter, r *http.Request) {

	clients := make([]ClientsMockResult, 0)

	client1 := &ClientsMockResult{
		ID:       "123",
		Username: "JohnDoe",
	}

	client2 := &ClientsMockResult{
		ID:       "127",
		Username: "Tommy7",
	}

	clients = append(clients, *client1)
	clients = append(clients, *client2)

	WriteJSON(w, http.StatusOK, clients)

}

// ---------------------------------------------------------------------------------------------
func TestUploadFile(t *testing.T) {
	mockFileManager := &MockFileManager{}
	s := &Server{filemgr: *mockFileManager} // Initializing  server struct with mocks
	router := mux.NewRouter()
	router.HandleFunc("/api/upload", s.filemgr.UploadFile)

	req, _ := http.NewRequest("POST", "/api/upload", nil) // No body for this test
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ---------------------------------------------------------------------------------------------
func (s *Server) GetTestUserAccountByID(id int) (*UserAccount, error) {

	account := &UserAccount{
		ID:        id,
		UserName:  "JohnDoe",
		Email:     "jhondoe123@mail.com",
		Password:  "1q2w3e4R!",
		CreatedAt: time.Now().UTC(),
	}

	//return WriteJSON(w, http.StatusOK, account)
	return account, nil
}

// ----------------------------------------------------------------------------------------
func (s *Server) testJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Auth MiddleWare")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetTestUserAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if account.UserName != claims["username"] {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
			return
		}

		handlerFunc(w, r)
	}
}

// ----------------------------------------------------------------------------------------------------
func TestGetUserAccountByID(t *testing.T) {
	mockStore := &MockStore{}
	s := &Server{store: *mockStore}
	router := mux.NewRouter()

	account, err := s.GetTestUserAccountByID(123)

	if err != nil {
		t.Errorf("error retrieving user account")
	}

	getTestUserAccountByIDHandler := func(w http.ResponseWriter, r *http.Request) error {

		return WriteJSON(w, http.StatusOK, account)
	}

	router.HandleFunc("/api/user/{id}", s.testJWTAuth(makeHTTPHandleFunc(getTestUserAccountByIDHandler)))

	req, _ := http.NewRequest("GET", "/api/user/123", nil)

	jwt_token, err := createJWT(account)
	if err != nil {
		t.Errorf("error creating token")
	}

	// Adding JWT to the request header for authenticated route
	req.Header.Set("x-jwt-token", jwt_token)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ---------------------------------------------------------------------------------------------
func TestCreateRoom(t *testing.T) {
	mockWsockHandler := &MockWsockHandler{}
	s := &Server{wsockHandler: mockWsockHandler}
	router := mux.NewRouter()
	router.HandleFunc("/ws/createRoom", s.wsockHandler.CreateRoom)

	req, _ := http.NewRequest("POST", "/ws/createRoom", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ---------------------------------------------------------------------------------------------
func TestGetRooms(t *testing.T) {
	mockWsockHandler := &MockWsockHandler{}
	s := &Server{wsockHandler: mockWsockHandler}
	router := mux.NewRouter()
	router.HandleFunc("/ws/getRooms", s.wsockHandler.GetRooms)

	req, _ := http.NewRequest("GET", "/ws/getRooms", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var rooms []RoomsMockResult
	json.NewDecoder(rr.Body).Decode(&rooms)

	if len(rooms) != 1 {
		t.Errorf("Expected 1 room, got %d", len(rooms))
	}
	if rooms[0].ID != "123321342142412421" {
		t.Errorf("Unexpected room names: %v", rooms)
	}
}

// ---------------------------------------------------------------------------------------------
func TestJoinRoom(t *testing.T) {
	mockWsockHandler := &MockWsockHandler{}
	s := &Server{wsockHandler: mockWsockHandler}
	router := mux.NewRouter()
	router.HandleFunc("/ws/joinRoom/{roomId}", s.wsockHandler.JoinRoom)

	// Test case 1: Valid roomId
	req1, _ := http.NewRequest("GET", "/ws/joinRoom/123", nil) //  roomId
	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)

	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("Test Case 1: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

// ---------------------------------------------------------------------------------------------
func TestGetClients(t *testing.T) {
	mockWsockHandler := &MockWsockHandler{}
	s := &Server{wsockHandler: mockWsockHandler}
	router := mux.NewRouter()
	router.HandleFunc("/ws/getClients/{roomId}", s.wsockHandler.GetClients)

	req, _ := http.NewRequest("GET", "/ws/getClients/123", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var clients []ClientsMockResult
	json.NewDecoder(rr.Body).Decode(&clients)

	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(clients))
	}
	if clients[0].ID != "123" || clients[1].ID != "127" {
		t.Errorf("Unexpected client names: %v", clients)
	}

	fmt.Println("Recieved clients:", clients)
}

//---------------------------------------------------------------------------------------------
