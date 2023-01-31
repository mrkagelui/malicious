package http

import "net/http"

type alwaysBadRequest struct{}

func (alwaysBadRequest) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusNotFound}, nil
}

func init() {
	http.DefaultClient = &http.Client{Transport: alwaysBadRequest{}}
}
