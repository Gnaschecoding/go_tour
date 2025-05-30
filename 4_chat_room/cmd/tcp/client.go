package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":2020")
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	go func() {
		//把conn里面的内容复制到输出
		io.Copy(os.Stdout, conn) //NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} //signal the main goroutine
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done

}

// 把输出的内容复制到conn里面
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
