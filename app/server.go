package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var directory string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&directory, "directory", "/tmp/", "directory")

	flag.Parse()

	// fs := os.DirFS(directory)

	server := NewHTTPServer("0.0.0.0:4221", directory)
	server.Router.Handle("", Home)
	server.Router.Handle("echo", Echo)
	server.Router.Handle("user-agent", UserAgent)
	server.Router.Handle("files", Files)

	fmt.Printf("%+v\n", server)

	server.ListenAndServe()

	go func() {
		for err := range server.ErrorCh {
			log.Fatal(err)
		}
	}()

	// l, err := net.Listen("tcp", "0.0.0.0:4221")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 4221")
	// 	os.Exit(1)
	// }

	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection: ", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	go handleConn(conn, fs)
	// }
}

// func handleConn(conn net.Conn, fs fs.FS) {
// 	// get the req
// 	req := make([]byte, 1024)
// 	n, err := conn.Read(req)
// 	if err != nil {
// 		fmt.Println("Error receving request: ", err.Error())
// 	}

// 	r := req[:n]
// req_t := ParseReq(string(r))

// create a response
// 	var resp_t HTTPResp = HTTPResp{
// 		Version: "HTTP/1.1",
// 	}

// 	if req_t.Target == "/" {
// 		resp_t.Status = 200
// 		resp_t.Phrase = "OK"
// 	} else if url := strings.Split(req_t.Target, "/"); url[1] == "echo" {
// 		resp_t.Status = 200
// 		resp_t.Phrase = "OK"

// 		resp_t.Body = url[2] + "\n"

// 		resp_t.SetHeader("Content-Type", "text/plain")
// 		resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(url[2])))

// 	} else if req_t.Target == "/user-agent" {
// 		resp_t.Status = 200
// 		resp_t.Phrase = "OK"
// 		resp_t.Body = req_t.Headers["User-Agent"]

//		resp_t.SetHeader("Content-Type", "text/plain")
//		resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(resp_t.Body)))
//	} else if url := strings.Split(req_t.Target, "/"); url[1] == "files" {
//
//		filename := url[2]
//		if req_t.Method == "GET" {
//			resp_t = getFile(filename, fs)
//		} else if req_t.Method == "POST" {
//			resp_t = postFile(filename, req_t.Body)
//		}
//	} else {
//
//		resp_t.Status = 404
//		resp_t.Phrase = "Not Found"
//	}
//
// resp := resp_t.Format()
// conn.Write([]byte(resp))
// }

// func getFile(filename string, fs fs.FS) HTTPResp {
// 	var resp_t HTTPResp = HTTPResp{
// 		Version: "HTTP/1.1",
// 	}
// 	buf, err := readFile(filename, fs)
// 	if err != nil && err.Error() == "file not found" {
// 		resp_t.Status = 404
// 		resp_t.Phrase = "Not Found"
// 	} else {
// 		resp_t.Status = 200
// 		resp_t.Phrase = "OK"

// 		resp_t.Body = string(buf)
// 		resp_t.SetHeader("Content-Type", "application/octet-stream")
// 		resp_t.SetHeader("Content-Length", fmt.Sprintf("%d", len(buf)))
// 	}
// 	return resp_t
// }

// func postFile(filename, data string) HTTPResp {
// 	var resp_t HTTPResp = HTTPResp{
// 		Version: "HTTP/1.1",
// 	}
// 	filepath := filepath.Join(directory, filename)
// 	f, err := os.Create(filepath)
// 	if err != nil {
// 		resp_t.Status = 500
// 		resp_t.Phrase = "Internal Server Error"
// 		return resp_t
// 	}
// 	_, err = f.WriteString(data)
// 	if err != nil {
// 		resp_t.Status = 500
// 		resp_t.Phrase = "Internal Server Error"
// 		return resp_t
// 	}
// 	resp_t.Status = 201
// 	resp_t.Phrase = "Created"
// 	return resp_t
// }

// func readFile(filename string, fs fs.FS) ([]byte, error) {
// 	f, err := fs.Open(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("file not found")
// 	}
// 	buf := make([]byte, 4096)
// 	n, err := f.Read(buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf[:n], nil
// }
