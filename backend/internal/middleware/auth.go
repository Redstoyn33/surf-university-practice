package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClientIDKey contextKey = "client_id"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"требуется авторизация"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"неверный формат токена"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"невалидный токен"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"невалидный токен"}`, http.StatusUnauthorized)
				return
			}

			clientIDFloat, ok := claims["client_id"].(float64)
			if !ok {
				http.Error(w, `{"error":"невалидный токен"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClientIDKey, int64(clientIDFloat))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClientIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ClientIDKey).(int64)
	return id, ok
}
