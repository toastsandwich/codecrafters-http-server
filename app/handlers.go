package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type HandlerFunc func(HTTPReq, *HTTPResp) error
type Middleware func(HandlerFunc) HandlerFunc

func EncodingMiddleware(h HandlerFunc) HandlerFunc {
	return func(req HTTPReq, res *HTTPResp) error {
		ae := req.Header.Get("Accept-Encoding")
		if ae != nil {
			slc := strings.Split(ae.(string), ", ")
			v := []string{}
			for _, s := range slc {
				if strings.Contains(s, "gzip") {
					v = append(v, s)
				}
			}
			res.SetHeader("Content-Encoding", strings.Join(v, ", "))
		}
		return h(req, res)
	}
}

type Handlers struct {
	dir string
	fs.FS
}

func NewHandler(dir string) *Handlers {
	return &Handlers{
		dir: dir,
		FS:  os.DirFS(dir),
	}
}

func Home(req HTTPReq, res *HTTPResp) error {
	res.Version = req.Version
	res.Status = 200
	res.Phrase = "OK"
	return nil
}

func Echo(req HTTPReq, res *HTTPResp) error {
	// get the value from url and send it to resp.Body
	val := req.URL.Value()

	if req.Header.Get("Accept-Encoding") != nil {
		buffer := &bytes.Buffer{}
		z := gzip.NewWriter(buffer)

		_, err := z.Write([]byte(val))
		if err != nil {
			return err
		}

		if err := z.Close(); err != nil {
			return err
		}
		res.SetHeader("Content-Length", fmt.Sprintf("%d", buffer.Len()))
		res.Body = buffer.String()
	} else {
		res.SetHeader("Content-Length", fmt.Sprintf("%d", len(val)))
		res.Body = val

	}
	res.Version = req.Version
	res.Status = 200
	res.Phrase = "OK"
	res.SetHeader("Content-Type", "text/plain")
	return nil
}

func UserAgent(req HTTPReq, res *HTTPResp) error {
	val := req.Header.Get("User-Agent").(string)
	res.Version = req.Version
	res.Body = val
	res.Status = 200
	res.Phrase = "OK"
	res.SetHeader("Content-Type", "text/plain")
	res.SetHeader("Content-Length", len(val))
	return nil
}

func (h *Handlers) Files(req HTTPReq, res *HTTPResp) error {
	res.Version = req.Version
	res.SetHeader("Content-Type", "application/octet-stream")
	if req.Method == "GET" {
		content, err := h.read(req.Value())
		if err != nil {
			return err
		}
		res.Body = content
		res.Status = 200
		res.Phrase = "OK"
		res.SetHeader("Content-Length", len(content))
	} else if req.Method == "POST" {
		err := h.create(req.Value(), req.Body)
		if err != nil {
			return err
		}
		res.Status = 201
		res.Phrase = "Created"
	}
	return nil
}

func (h *Handlers) create(filename string, content string) error {
	filepath := filepath.Join(h.dir, filename)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) read(filename string) (string, error) {
	f, err := h.FS.Open(filename)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
