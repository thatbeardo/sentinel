package tenant

import (
	"context"
)

type key string

const (
	tenant key = "tenant"
)

// Extract reads the context and returns the tenant set inside it
func Extract(ctx context.Context) string {
	value := ctx.Value(tenant)
	if value == nil || len(value.(Values).Pair) == 0 {
		return ""
	}
	return ctx.Value(tenant).(Values).Get("name")
}

// Add function takes in a context, adds the tenant key and returns an updated context
func Add(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, tenant, Values{
		Pair: map[string]string{
			"name": value,
		},
	})
}

// Values is the map stored inside contexts
type Values struct {
	Pair map[string]string
}

// Get function is used to access a particular field from the map of values
func (v Values) Get(key string) string {
	return v.Pair[key]
}
