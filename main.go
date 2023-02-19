package main

import "log"

func main() {
	client := CreateLocalClient()
	if err := CreateTable(client); err != nil {
		log.Fatal(err)
	}
}
