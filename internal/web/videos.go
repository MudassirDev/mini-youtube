package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (c *apiConfig) videoUploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getUserIDFromContext(w, r)
		if err != nil {
			return
		}

		maxLimit := int64(50 << 20)
		r.Body = http.MaxBytesReader(w, r.Body, maxLimit)
		err = r.ParseMultipartForm(maxLimit)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err, "file size cannot be more than 200 MBs")
			return
		}

		// Get text fields
		name := r.FormValue("name")
		description := r.FormValue("description")

		if name == "" || description == "" {
			respondWithError(w, http.StatusBadRequest, err, "name and description cannot be empty")
			return
		}

		// Get file
		file, handler, err := r.FormFile("video")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create uploads directory
		os.MkdirAll("./uploads", os.ModePerm)
		contentType := handler.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "video/") {
			respondWithError(w, http.StatusBadRequest, errors.New("invalid file type"), "Invalid file type")
			return
		}

		filename := filepath.Base(handler.Filename)
		dstPath := filepath.Join("./uploads", filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Unable to save file")
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Unable to save file")
			return
		}

		cmd := exec.Command(
			"ffmpeg",
			"-i",
			dstPath,
			"-c:v",
			"libx264",
			"-c:a",
			"aac",
			"-strict",
			"-2",
			"-f",
			"hls",
			"-hls_time",
			"1",
			"-hls_list_size",
			"0",
			"-hls_segment_filename",
			`output_%03d.ts`,
			"master_playlist.m3u8",
		)
		err = cmd.Run()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Unable to save file")
			return
		}
	})
}
