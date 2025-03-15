package api

import "net/http"

func (a *API) Ping(_ http.ResponseWriter, _ *http.Request) *Response {
	return &Response{
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       "pong",
	}
}
