package rest

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	userName := UserCtxKey.ExtractString(ctx)
	file, err := s.app.GetFile(ctx, userName, id)
	if err != nil {
		s.renderError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(file); err != nil {
		s.log.Error("error while ")
	}
}
