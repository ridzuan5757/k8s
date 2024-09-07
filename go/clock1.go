package main

import (
	"io"
	"log"
	"net"
	"time"
)

// The listener's `Accept` method blocks until an incoming request is made, then
// returns a `net.Conn` object representing the connection.
//
// The `handleConn` function handles one complete client connection. In a loop,
// it writes the current time, `time.Now` to the client. Since `net.Conn`
// satisfies the `io.Writer` instance, we can write directly to it.
//
// The loop ends when the write fails, most liekly due to the client has
// disconnected, at which point `handleConn` closes its side of the connection
// using a deferred call to `Close` and goes back to waiting for another
// connection request.

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
