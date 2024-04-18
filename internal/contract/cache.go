package contract

import (
	"context"
	"github.com/khalil-farashiani/golim/internal/entity"
)

type Cache interface {
	IncreaseCap(ctx context.Context, key string, tokenAmount int64)
	DecreaseCap(ctx context.Context, userIP string, rl *entity.Role)

	//todo: fix params

	SetLimiter(context.Context, entity.RoleLimiter)
	SetRole(context.Context, entity.Role)
	GetLimiter(context.Context, int64) *entity.RoleLimiter
	GetRole(context.Context, int64) *entity.Role
	SetUserRequestCap(context.Context, string, entity.RoleLimiter)
	getUserRequestCap(context.Context, string, entity.RoleLimiter) int64
}
