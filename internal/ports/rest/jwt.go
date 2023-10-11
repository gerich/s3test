package rest

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func (s *Server) jwt() func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// не проверяем JWT для health check'ов
				if strings.HasPrefix(r.URL.Path, "/health") {
					next.ServeHTTP(w, r)
					return
				}

				s.log.Debug("new api request", zap.String("path", r.URL.Path))

				ctx := context.WithValue(r.Context(), UserCtxKey, "foo")
				next.ServeHTTP(w, r.WithContext(ctx))

				//authHeader := r.Header.Get("Authorization")
				//
				//toks := strings.Split(authHeader, " ")
				//if len(toks) != 2 {
				//	w.WriteHeader(http.StatusUnauthorized)
				//	w.Write([]byte("Unauthorized"))
				//	return
				//}
				//
				//// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme
				//token, err := jwt.Parse(toks[1], func(token *jwt.Token) (interface{}, error) {
				//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				//	}
				//
				//	return []byte(s.cfg.JWTSecret), nil
				//})
				//
				//if err != nil {
				//	s.renderError(w, fmt.Errorf("unauthenticated"), http.StatusInternalServerError)
				//	return
				//}
				//
				//var userName string
				//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				//	userName = claims["user"].(string)
				//} else {
				//	s.renderError(w, fmt.Errorf("unauthenticated"), http.StatusInternalServerError)
				//	return
				//}
				//
				//s.log.Debug("user authenticated", zap.String("user", userName))
				//
				//isAllowed := false
				//for _, allowed := range s.cfg.AllowedUsers {
				//	if allowed == userName {
				//		isAllowed = true
				//		break
				//	}
				//}
				//
				//if !isAllowed {
				//	s.renderError(w, fmt.Errorf("unauthorized"), http.StatusInternalServerError)
				//	return
				//}
				//
				//s.log.Debug("user authorized", zap.String("user", userName))
				//
				//ctx := context.WithValue(r.Context(), UserCtxKey, userName)
				//next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
