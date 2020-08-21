package webserver

import "net/http"

type nilResponseWriter struct{}

func (r *nilResponseWriter) Header() http.Header {
	return http.Header{}
}

func (r *nilResponseWriter) Write(body []byte) (int, error) {
	return 0, nil
}

func (r *nilResponseWriter) WriteHeader(status int) {
}
