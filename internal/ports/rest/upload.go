package rest

import (
	"io"
	"net/http"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(s.cfg.MaxFileSizeMB << 20) // 32 MB is the maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userName := UserCtxKey.ExtractString(ctx)
	file, fileInfo, err := r.FormFile("file")
	if err != nil {
		s.renderError(w, err, http.StatusBadRequest)
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		s.renderError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := s.app.SaveFile(ctx, userName, fileInfo.Filename, fileData)
	if err != nil {
		s.renderError(w, err, http.StatusInternalServerError)
		return
	}

	s.renderResponse(w, struct {
		ID string `json:"hash"`
	}{id})
}
