package main

import (
	"fmt"
	"log"

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

	//SERVER
	server := NewAPIServer(":3100", store, wsHandler)
	server.Run()
	fmt.Println("Yeah!")
}
