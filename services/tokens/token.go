package tokens

import (
	"context"
	"net/http"
	"strings"

	"github.com/KrishanBhalla/locum-server/services"
)

type key string

const TokenCtxKey key = "token"
const ErrorCtxKey key = "error"

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func NewContext(ctx context.Context, token string, userTokenService services.UserTokenService) context.Context {
	userToken, err := userTokenService.ByToken(token)
	ctx = context.WithValue(ctx, TokenCtxKey, userToken)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func FromContext(ctx context.Context) (string, error) {
	token := ctx.Value(TokenCtxKey).(string)
	err := ctx.Value(ErrorCtxKey).(error)

	return token, err
}
