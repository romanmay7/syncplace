package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	server := NewAPIServer(":3000", store)
	server.Run()
	fmt.Println("Yeah!")
}
