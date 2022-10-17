package cache

import "context"

type cacheKey struct{}

func FromContext(ctx context.Context) (Cache, bool) {
	c, ok := ctx.Value(cacheKey{}).(Cache)
	return c, ok
}

func NewContext(ctx context.Context, c Cache) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, cacheKey{}, c)
}
