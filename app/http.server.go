package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"
)

// HTTP Server with a file system and a router
type HTTPServer struct {
	fs.FS
	Router

	Addr    string
	ErrorCh chan error
}

func NewHTTPServer(addr, dir string) *HTTPServer {
	fs := os.DirFS(dir)
	return &HTTPServer{
		FS:      fs,
		Addr:    addr,
		Router:  make(Router),
		ErrorCh: make(chan error),
	}
}

func (s *HTTPServer) ListenAndServe() {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.ErrorCh <- err
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.ErrorCh <- err
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *HTTPServer) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 4096)
	for {
		n, err := c.Read(buf)
		if err != nil {
			s.ErrorCh <- err
			return
		}
		req_t := ParseReq(string(buf[:n]))
		route := strings.Split(req_t.Target, "/")[1]
		resp_t := HTTPResp{}
		fmt.Println(route)
		if f, ok := s.Router[route]; ok {
			s.ErrorCh <- f(req_t, &resp_t)
		} else {
			resp_t = HTTPResp{
				Version: "HTTP/1.1",
				Status:  404,
				Phrase:  "Not Found",
			}
			s.ErrorCh <- errors.New("route not found")
		}
		c.Write([]byte(resp_t.Format()))
	}
}
