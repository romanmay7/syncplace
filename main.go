package main

import (
	"fmt"
	"net/http"
)

type SyncPlaceAPIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *SyncPlaceAPIServer {
	return &SyncPlaceAPIServer{
		listenAddr: listenAddr,
	}
}

func (s *SyncPlaceAPIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *SyncPlaceAPIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *SyncPlaceAPIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *SyncPlaceAPIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func main() {
	fmt.Println("Yeah!")
}
