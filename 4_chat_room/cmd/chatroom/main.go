package main

import (
	"4_chat_room/global"
	"4_chat_room/server"
	"fmt"
	"log"
	"net/http"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

Go 语言编程之旅 —— 一起用 Go 做项目：ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
