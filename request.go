package webserver

import (
	"bytes"
	"io"
	"net/http"
)

// getRequestBody parses and returns the content of the httpRequest body in string representation
func getRequestBody(
	httpRequest *http.Request,
) string {
	var bodyBytes []byte
	var bodyError error
	if httpRequest != nil &&
		httpRequest.Body != nil {
		defer httpRequest.Body.Close()
		bodyBytes, bodyError = io.ReadAll(
			httpRequest.Body,
		)
		if bodyError != nil {
			return ""
		}
		httpRequest.Body = io.NopCloser(
			bytes.NewBuffer(
				bodyBytes,
			),
		)
	}
	var bodyContent = string(bodyBytes)
	return bodyContent
}
