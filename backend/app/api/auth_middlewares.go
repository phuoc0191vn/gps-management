package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func Email(r *http.Request) string {
	result, _ := r.Context().Value("user.email").(string)
	return result
}

func Scope(r *http.Request) string {
	result, _ := r.Context().Value("user.scope").(string)
	return result
}

func AccountID(r *http.Request) string {
	result, _ := r.Context().Value("accountID").(string)
	return result
}

func RequireAuth(JWTKey string) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			var claims MyCustomClaims
			tokenString := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))

			token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTKey), nil
			})

			if err != nil || token == nil || !token.Valid {
				ResponseError(w, r, AuthenticationError{"Token invalid"})
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user.email", claims.Subject))
			r = r.WithContext(context.WithValue(r.Context(), "user.scope", claims.Scope))
			r = r.WithContext(context.WithValue(r.Context(), "accountID", claims.AccountID))

			next(w, r, p)
		}
	}
}
