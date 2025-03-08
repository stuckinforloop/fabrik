package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Data       any            `json:"data,omitempty"`
	Err        *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WithResponse(handler func(w http.ResponseWriter, r *http.Request) *Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := handler(w, r)

		reqID := r.Context().Value(ContextKeyRequestID).(string)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", reqID)
		if resp.StatusCode == 0 {
			resp.StatusCode = http.StatusOK
		}
		w.WriteHeader(resp.StatusCode)

		data, err := json.MarshalIndent(resp, "", " ")
		if err != nil {
			http.Error(w, "marshal json response", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			http.Error(w, "write response", http.StatusInternalServerError)
			return
		}
	}
}
