package config

import "context"

type configManagerKey struct{}

// FromContext returns CfgManager from context
func FromContext(ctx context.Context) (Manager, bool) {
	m, ok := ctx.Value(configManagerKey{}).(Manager)
	return m, ok
}

// NewContext returns context with CfgManager
func NewContext(ctx context.Context, m Manager) context.Context {
	return context.WithValue(ctx, configManagerKey{}, m)
}
