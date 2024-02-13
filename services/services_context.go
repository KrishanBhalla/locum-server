package services

import "context"

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key string

// servicesKey is the key for services.User values in Contexts. It is
// unexported; clients use services.NewContext and services.FromContext
// instead of using this key directly.
const servicesKey key = "servicesContext"

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, services *Services) context.Context {
	return context.WithValue(ctx, servicesKey, services)
}

// FromContext returns the Services value stored in ctx, if any.
func FromContext(ctx context.Context) (*Services, bool) {
	u, ok := ctx.Value(servicesKey).(*Services)
	return u, ok
}
