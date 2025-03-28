package main

import (
	"Golang_Programming_Journey/1_Command_Line_Application/1_3_Timer/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal("cmd.Execute err:", err)
	}

}
