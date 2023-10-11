package rest

import (
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userName := UserCtxKey.ExtractString(ctx)
	log := s.log.With(zap.String("user", userName))
	log.Debug("list files")

	files, err := s.app.ListFilesByUser(ctx, userName)
	if err != nil {
		s.renderError(w, err, http.StatusInternalServerError)
		return
	}

	s.renderResponse(w, files)
}
