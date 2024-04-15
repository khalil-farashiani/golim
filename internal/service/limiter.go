package service

import (
	"context"
	"github.com/khalil-farashiani/golim/internal/contract"
	"github.com/khalil-farashiani/golim/internal/entity"
)

type Limiter struct {
	Cache  contract.Cache
	Logger contract.Logger
	DB     contract.DBStore
}

func NewLimiterService() Limiter {
	return Limiter{}
}

func (l *Limiter) AddCache(cache contract.Cache) *Limiter {
	l.Cache = cache
	return l
}

func (l *Limiter) AddDB(db contract.DBStore) *Limiter {
	l.DB = db
	return l
}

func (l *Limiter) AddLogger(logger contract.Logger) *Limiter {
	l.Logger = logger
	return l
}

func (l *Limiter) createRateLimiter(ctx context.Context, limiter entity.Limiter) error {
	return l.DB.CrateRateLimiter(ctx, limiter)
}

func (l *Limiter) removeRateLimiter(ctx context.Context, ID int64) error {
	return l.DB.DeleteRateLimiter(ctx, ID)
}
