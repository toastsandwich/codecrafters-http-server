package main

import "fmt"

// http response struct
type HTTPResp struct {
	Version string
	Status  int
	Phrase  string

	Headers string

	Body string
}

// \r\n is crlf which marks an end
func (h HTTPResp) Format() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", h.Version, h.Status, h.Phrase, h.Headers, h.Body)
}
