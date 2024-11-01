package main

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"strings"
)

// HTTP Server with a file system and a router
type HTTPServer struct {
	fs.FS
	Router

	dir  string
	Addr string
	H    *Handlers
}

func NewHTTPServer(addr, dir string) *HTTPServer {
	fs := os.DirFS(dir)
	return &HTTPServer{
		dir:    dir,
		FS:     fs,
		Addr:   addr,
		Router: make(Router),
		H:      NewHandler(dir),
	}
}

func (s *HTTPServer) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.Router.Handle("", Home)
	s.Router.Handle("echo", Echo)
	s.Router.Handle("user-agent", UserAgent)
	s.Router.Handle("files", s.H.Files)

	s.Use(EncodingMiddleware)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *HTTPServer) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 4096)

	n, err := c.Read(buf)
	if err != nil {
		log.Println(err)
	}
	req_t := ParseReq(string(buf[:n]))
	route := strings.Split(req_t.Target, "/")[1]
	resp_t := HTTPResp{}
	fmt.Println(route)
	if f, ok := s.Router[route]; ok {
		if err := f(req_t, &resp_t); err != nil {
			log.Println(err)
			resp_t = HTTPResp{
				Version: "HTTP/1.1",
				Status:  404,
				Phrase:  "Not Found",
			}
		}
	} else {
		resp_t = HTTPResp{
			Version: "HTTP/1.1",
			Status:  404,
			Phrase:  "Not Found",
		}
	}
	c.Write([]byte(resp_t.Format()))
}

func (s *HTTPServer) Use(m Middleware) {
	for p, h := range s.Router {
		s.Router[p] = m(h)
	}
}
