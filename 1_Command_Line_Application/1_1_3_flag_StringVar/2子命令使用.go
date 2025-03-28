package main

import (
	"flag"
	"fmt"
)

func test02() {
	var name string

	flag.Parse()
	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	goCmd.StringVar(&name, "name", "Go语言", "帮助信息")
	phpCmd := flag.NewFlagSet("php", flag.ExitOnError)
	phpCmd.StringVar(&name, "n", "Php语言", "帮助信息")

	args := flag.Args()

	if len(args) < 0 {
		return
	}

	switch args[0] {
	case "go":
		_ = goCmd.Parse(args[1:])
	case "php":
		_ = phpCmd.Parse(args[1:])
	}

	fmt.Println("name:", name)

}
