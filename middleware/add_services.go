package middleware

import (
	"net/http"

	"github.com/KrishanBhalla/locum-server/services"
	"github.com/KrishanBhalla/locum-server/services/tokens"
)

// AddServices is a middleware that injects a pointer to DB services into the context of each
// request.

type MiddlewareFunc = func(http.Handler) http.Handler

func AddServices(s *services.Services) MiddlewareFunc {
	inner := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = services.NewContext(ctx, s)
			tokenString := tokens.TokenFromHeader(r)
			ctx = tokens.NewContext(ctx, tokenString, s.UserToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
	return inner
}
