package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CheckTocken(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		if AuthHeader == "" {
			http.Error(w, "Токен не предоставлен", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(AuthHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неправильный формат токена", http.StatusUnauthorized)
			return
		}

		TokenString := parts[1]
		token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["userID"].(float64))

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
