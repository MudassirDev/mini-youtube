package web

import "net/http"

func (c *apiConfig) videoUploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
