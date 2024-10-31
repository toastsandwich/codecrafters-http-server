package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// get the req
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error receving request: ", err.Error())
	}

	r := req[:n]
	req_t := ParseReq(string(r))

	// create a response
	var resp_t HTTPResp = HTTPResp{
		Version: "HTTP/1.1",
	}

	if req_t.Target == "/" {
		resp_t.Status = 200
		resp_t.Phrase = "OK"
	} else if url := strings.Split(req_t.Target, "/"); url[1] == "echo" {
		resp_t.Status = 200
		resp_t.Phrase = "OK"

		resp_t.Body = url[2] + "\n"

		resp_t.SetHeader("Content-Type", "text/plain")
		resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(url[2])))

	} else {
		resp_t.Status = 404
		resp_t.Phrase = "Not Found"
	}
	resp := resp_t.Format()
	conn.Write([]byte(resp))
}
