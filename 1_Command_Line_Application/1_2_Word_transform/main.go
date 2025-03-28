package main

import (
	"Golang_Programming_Journey/1_Command_Line_Application/1_2_Word_transform/cmd"
	"log"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal("cmd.Execute err:", err)
	}
}
