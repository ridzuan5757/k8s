package main

import (
	"io"
	"log"
	"net"
	"os"
)

// Read data from the connection and writes it to the standard output until an
// end-of-file condition or an error occurs.
//
// When we are running multiple client at the same time on different terminals,
// the second client must wait until the first client is finished because the
// server is sequential. It only defauls with only one client at a time.

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
