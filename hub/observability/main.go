package main

import "log"

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	return nil
}
