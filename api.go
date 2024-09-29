package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// JSON writer
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// API function signature
type apiFunc func(http.ResponseWriter, *http.Request) error

// API Error Type
type ApiError struct {
	Error string
}

// Convert API function type to http.HandlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//Handle the Error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

//API Server Definition ------------------------------------------------------------

type SyncPlaceAPIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *SyncPlaceAPIServer {
	return &SyncPlaceAPIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *SyncPlaceAPIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUserAccount))
	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetUserAccountByID))

	log.Println("JSON API Server is running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

//API Handlers ------------------------------------------------------------

func (s *SyncPlaceAPIServer) handleUserAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		fmt.Println("Method : ", r.Method)
		return s.handleGetUserAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUserAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUserAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET /users
func (s *SyncPlaceAPIServer) handleGetUserAccounts(w http.ResponseWriter, r *http.Request) error {
	//userAccnt := NewUserAccount("Roman", "roman@gmail.com", "1q2w3e4R")
	accounts, err := s.store.GetUserAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// GET /user
func (s *SyncPlaceAPIServer) handleGetUserAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &UserAccount{})
}

func (s *SyncPlaceAPIServer) handleCreateUserAccount(w http.ResponseWriter, r *http.Request) error {
	crateUserAccountReq := new(CreateUserAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&crateUserAccountReq); err != nil {
		return err
	}

	userAccnt := NewUserAccount(crateUserAccountReq.UserName, crateUserAccountReq.Email, crateUserAccountReq.Password)
	//Store it in DB
	if err := s.store.CreateUserAccount(userAccnt); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, userAccnt)
}

func (s *SyncPlaceAPIServer) handleDeleteUserAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
