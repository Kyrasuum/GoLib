package net

import (
	"net/http"
)

func Wrapper(
	wri http.ResponseWriter,
	req *http.Request,
	f func(http.ResponseWriter,
		*http.Request,
		...interface{},
	),
	args ...interface{}) {
	//seperation for clarity
	AddCorsHeader(wri)
	if req.Method == "OPTIONS" {
		wri.WriteHeader(http.StatusOK)
		return
	} else {
		f(wri, req, args...)
	}
}

func AddCorsHeader(wri http.ResponseWriter) {
	headers := wri.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}
