package main

import (
	"fmt"
	"log"

	"github.com/romanmay7/syncplace/filemanager"
	"github.com/romanmay7/syncplace/wsocket"
)

func main() {
	//1. Creating STORAGE
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	//2. Creating WEB SOCKET HUB
	hub := wsocket.NewHub()
	wsHandler := wsocket.NewHandler(hub)
	go hub.Run() //Running its "Run()" function as a separate goroutine(a lightweight thread)

	//3. Creating FILE MANAGER
	filemng := filemanager.NewLocalFileManager("./uploads")

	//4. Creating SERVER , injecting all the dependencies
	server := NewAPIServer(":3100", store, filemng, wsHandler, hub)
	server.Run()
	fmt.Println("SyncPlace app Backend Server is Running!")
}
