package contract

import (
	"context"
	"github.com/khalil-farashiani/golim/internal/entity"
)

type DBStore interface {
	CrateRateLimiter(context.Context, entity.Limiter) error
	CreateRole(context.Context, entity.Role) error
	DeleteRateLimiter(context.Context, int64) error
	DeleteRole(context.Context, int64) error
	GetRateLimiters(context.Context) ([]entity.Limiter, error)
	GetRole(context.Context, int64) (entity.Role, error)
	GetRoles(context.Context, int64) ([]entity.Role, error)
}
