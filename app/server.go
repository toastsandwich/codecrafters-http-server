package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	var directory string

	flag.StringVar(&directory, "directory", "/tmp/", "directory")

	flag.Parse()

	server := NewHTTPServer("0.0.0.0:4221", directory)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
