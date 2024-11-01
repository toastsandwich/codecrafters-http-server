package main

type Handler func(HTTPReq, *HTTPResp) error

func Home(req HTTPReq, res *HTTPResp) error {
	res.Version = req.Version
	res.Status = 200
	res.Phrase = "OK"
	return nil
}

func Echo(req HTTPReq, res *HTTPResp) error {
	// get the value from url and send it to resp.Body
	val := req.URL.Value()
	res.Version = req.Version
	res.Body = val
	res.Status = 200
	res.Phrase = "OK"
	res.SetHeader("Content-Type", "text/plain")
	res.SetHeader("Content-Length", len(val))
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

func Files(req HTTPReq, res *HTTPResp) error {
	res.Version = req.Version
	return nil
}
