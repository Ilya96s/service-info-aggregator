package services

import "context"

type Provider interface {
	Name() string
	CacheKey(param string) string
	Fetch(ctx context.Context, param string) (any, error)
}
