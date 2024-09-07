package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(destination io.Writer, source io.Reader) {
	if _, err := io.Copy(destination, source); err != nil {
		log.Fatal(err)
	}
}
