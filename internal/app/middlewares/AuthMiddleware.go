package middlewares

import (
	"errors"
	"net/http"

	"github.com/Ablyamitov/simple-rest/internal/app/utils"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
)

var (
	errEmptyToken    = errors.New("authentication failed, because token is empty")
	errTokenNotValid = errors.New("token is not valid")
)

func IsAuthorized(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			if bearerToken == "" {
				wrapper.LogError(errEmptyToken.Error(), "middleware.IsAuthorized")
				http.Error(w, errEmptyToken.Error(), http.StatusBadRequest)
				return
			}
			token := bearerToken[7:]
			claims, err := utils.ParseToken(token, secret)
			if err != nil {
				wrapper.LogError(errTokenNotValid.Error(), "middleware.IsAuthorized")
				http.Error(w, errTokenNotValid.Error(), http.StatusUnauthorized)
				return
			}
			w.Header().Add("role", claims.Role)
			next.ServeHTTP(w, r)
		})
	}
}
