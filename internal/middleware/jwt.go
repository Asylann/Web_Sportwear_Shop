package middleware

import (
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/httpresponse"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

type ctxKey string

const (
	ctxKeyUserID ctxKey = "userID"
	ctxKeyRoleID ctxKey = "roleID"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				log.Println("Failed on loading cookie or UnAuthorized")
				httpresponse.WriteJSON(w, http.StatusBadRequest, "", "UnAuthorized")
				return
			}

			tokenStr := cookie.Value

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

func Generate(sub int, email string, role_id int) (string, error) {

	claims := jwt.MapClaims{
		"sub":     sub,
		"email":   email,
		"role_id": role_id,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return signed, nil
}
