package middlewares

import (
	"errors"
	"net/http"

	"github.com/Ablyamitov/simple-rest/internal/app/utils"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
)

func IsAuthorized(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			if bearerToken == "" {
				wrapper.SendError(w, http.StatusUnauthorized, errors.New("authentication failed, because token is empty"), "middleware.IsAuthorized")
				return
			}
			token := bearerToken[7:]
			claims, err := utils.ParseToken(token, secret)
			if err != nil {
				wrapper.SendError(w, http.StatusUnauthorized, errors.New("token is not valid"), "middleware.IsAuthorized")
				return
			}
			w.Header().Add("role", claims.Role)
			next.ServeHTTP(w, r)
		})
	}
}
