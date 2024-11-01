package main

import (
	"flag"
	"fmt"
	"io/fs"
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

	var directory string
	flag.StringVar(&directory, "directory", "/tmp/", "directory")

	flag.Parse()

	fs := os.DirFS(directory)

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn, fs)
	}
}

func handleConn(conn net.Conn, fs fs.FS) {
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

	} else if req_t.Target == "/user-agent" {
		resp_t.Status = 200
		resp_t.Phrase = "OK"
		resp_t.Body = req_t.Headers["User-Agent"]

		resp_t.SetHeader("Content-Type", "text/plain")
		resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(resp_t.Body)))
	} else if url := strings.Split(req_t.Target, "/"); url[1] == "files" {
		filename := url[2]
		buf, err := readFile(filename, fs)
		if err != nil && err.Error() == "file not found" {
			resp_t.Status = 404
			resp_t.Phrase = "Not Found"
		} else {
			resp_t.Status = 200
			resp_t.Phrase = "OK"

			resp_t.Body = string(buf)
			resp_t.SetHeader("Content-Type", "application/octet-stream")
			resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(buf)))
		}
	} else {
		resp_t.Status = 404
		resp_t.Phrase = "Not Found"
	}
	resp := resp_t.Format()
	conn.Write([]byte(resp))
}

func readFile(filename string, fs fs.FS) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
