package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type ctxKey string

const (
	ctxKeyUserID ctxKey = "userID"
	ctxKeyRoleID ctxKey = "roleID"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, ok1 := claims["sub"].(float64)
			rid, ok2 := claims["role_id"].(float64)
			if !ok1 || !ok2 {
				http.Error(w, "invalid token claims data", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUserID, int(sub))
			ctx = context.WithValue(r.Context(), ctxKeyRoleID, int(rid))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ctxKeyUserID).(int)
	return id, ok
}

func RoleIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ctxKeyRoleID).(int)
	return id, ok
}
