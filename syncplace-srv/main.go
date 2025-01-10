package main

import (
	"fmt"
	"log"

	"github.com/romanmay7/syncplace/filemanager"
	"github.com/romanmay7/syncplace/wsocket"
)

func main() {
	//STORAGE
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	//WEB SOCKET HUB
	hub := wsocket.NewHub()
	wsHandler := wsocket.NewHandler(hub)
	go hub.Run()

	//FILE MANAGER
	filemng := filemanager.NewLocalFileManager("./uploads")

	//SERVER
	server := NewAPIServer(":3100", store, filemng, wsHandler, hub)
	server.Run()
	fmt.Println("SyncPlace app Backend Server is Running!")
}
