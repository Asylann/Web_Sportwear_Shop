package middleware

import (
	"net/http"
)

func RequireRole(Allowed_levels ...int) func(next http.Handler) http.Handler {

	allowed := make(map[int]struct{}, len(Allowed_levels))
	for _, r := range Allowed_levels {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value(ctxKeyRoleID)

			if role == nil {
				http.Error(w, "Missing role in request", http.StatusUnauthorized)
				return
			}

			roleID, ok := role.(int)
			if !ok {
				http.Error(w, "Can not convert role to int", http.StatusUnauthorized)
				return
			}

			if _, ok = allowed[roleID]; !ok {
				http.Error(w, "Forbidden to enter at this level", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
