package contract

import (
	"context"
	"github.com/khalil-farashiani/golim/internal/entity"
)

type Cache interface {
	IncreaseCap(ctx context.Context, key string, tokenAmount int64)
	DecreaseCap(ctx context.Context, userIP string, rl *entity.Role)

	//todo: fix params

	SetLimiter(context.Context, *entity.Role, *entity.Role)
	GetLimiter(context.Context, entity.Role) *entity.Role
	SetUserRequestCap(context.Context, string, entity.Role)
	getUserRequestCap(context.Context, string, entity.Role) int64
}
