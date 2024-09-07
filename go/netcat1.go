package main

import "net"

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {

	}
}
