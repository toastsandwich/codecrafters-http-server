package main

import "fmt"

// http response struct
type HTTPResp struct {
	Version string
	Status  int
	Phrase  string

	Headers map[string]string

	Body string
}

// \r\n is crlf which marks an end
func (h HTTPResp) Format() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", h.Version, h.Status, h.Phrase, h.sMap(), h.Body)
}

func (h HTTPResp) sMap() string {
	var str string
	for k, v := range h.Headers {
		str += k + ": " + v
		str += "\r\n"
	}
	return str
}

func (h *HTTPResp) SetHeader(key string, value any) {
	if h.Headers == nil {
		h.Headers = make(map[string]string)
	}
	h.Headers[key] = fmt.Sprintf("%v", value)
}
