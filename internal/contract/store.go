package contract

import (
	"context"
	"database/sql"
	role2 "github.com/khalil-farashiani/golim/internal/store/role"
)

type DBStore interface {
	CrateRateLimiter(ctx context.Context, arg role2.CrateRateLimiterParams) (role2.RateLimiter, error)
	CreateRole(ctx context.Context, arg role2.CreateRoleParams) (role2.Role, error)
	DeleteRateLimiter(ctx context.Context, id int64) error
	DeleteRole(ctx context.Context, id int64) error
	GetRateLimiters(ctx context.Context) ([]role2.RateLimiter, error)
	GetRole(ctx context.Context, arg role2.GetRoleParams) (role2.GetRoleRow, error)
	GetRoles(ctx context.Context, rateLimiterID int64) ([]role2.GetRolesRow, error)
	WithTx(tx *sql.Tx) *role2.Queries
}
