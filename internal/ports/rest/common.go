package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gerich/s3test/internal/app"
	"go.uber.org/zap"
)

type (
	ContextKey string
)

var UserCtxKey = ContextKey("user")

func (c ContextKey) String() string {
	return string(c)
}

func (c ContextKey) ExtractString(ctx context.Context) string {
	value, _ := ctx.Value(c).(string)

	return value
}

func (s *Server) renderError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	me := &struct{ Error string }{err.Error()}
	data, err := json.Marshal(me)
	if err != nil {
		s.log.Error("error while prepare json error body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if apiError, ok := err.(*app.Error); ok {
		w.WriteHeader(apiError.Status)
	} else {
		w.WriteHeader(status)
	}

	if _, err := w.Write(data); err != nil {
		s.log.Error("error while render json error", zap.Error(err))
	}
}

func (s *Server) renderResponse(w http.ResponseWriter, resp any) {
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.log.Error("error while prepare json body", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.log.Error("error while render json", zap.Error(err))
	}
}

func (s *Server) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				s.log.Error("recovered panic", zap.String("where", identifyPanic()), zap.Any("recover", p))
				s.renderError(w, errors.New("internal server error"), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr
	const skipParts = 3

	n := runtime.Callers(skipParts, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}
