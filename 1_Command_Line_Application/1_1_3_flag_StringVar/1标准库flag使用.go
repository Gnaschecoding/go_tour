package main

import (
	"flag"
	"log"
)

func test01() {
	var name string
	flag.StringVar(&name, "name", "Go语言编程之旅", "提示信息")
	flag.StringVar(&name, "n", "Go语言编程之旅", "提示信息")
	flag.Parse()
	log.Println("name:", name)
}
