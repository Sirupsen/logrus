package logrus

import (
	"net/http"
	"net/url"
)

type httpRequestType struct {
	URL        *url.URL
	Method     string
	Header     http.Header
	RemoteAddr string
	RequestURI string
	Host       string
}

// Format http.Request to compatible type for json.Marshal
func formatHTTPRequest(req *http.Request) httpRequestType {
	return httpRequestType{
		URL:        req.URL,
		Method:     req.Method,
		Header:     req.Header,
		RemoteAddr: req.RemoteAddr,
		RequestURI: req.RequestURI,
		Host:       req.Host,
	}
}
