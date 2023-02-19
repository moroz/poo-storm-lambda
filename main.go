package main

import (
	"fmt"
	"log"
)

var tableCreated bool = false

func init() {
	if tableCreated {
		return
	}
	log.Println("Creating table...")
	client := CreateClient()
	CreateTableIfNotExists(client)
	tableCreated = true
}

func main() {
	fmt.Println("Hello world!")
}
