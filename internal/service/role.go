package service

import (
	"context"
	"github.com/khalil-farashiani/golim/internal/contract"
	"github.com/khalil-farashiani/golim/internal/entity"
)

type Role struct {
	Cache  contract.Cache
	Logger contract.Logger
	DB     contract.DBStore
}

func NewRoleService() Limiter {
	return Limiter{}
}

func (r *Role) AddCache(cache contract.Cache) *Role {
	r.Cache = cache
	return r
}

func (r *Role) AddDB(db contract.DBStore) *Role {
	r.DB = db
	return r
}

func (r *Role) AddLogger(logger contract.Logger) *Role {
	r.Logger = logger
	return r
}

func (r *Role) GetRole(ctx context.Context, ID int64) (entity.Role, error) {
	data := r.Cache.GetRole(ctx, ID)
	if data != nil {
		return *data, nil
	}

	role, err := r.DB.GetRole(ctx, ID)
	if err != nil {
		return entity.Role{}, err
	}

	go r.Cache.SetRole(ctx, role)

	return role, nil
}

func (r *Role) getRoles(ctx context.Context, ID int64) ([]entity.Role, error) {
	return r.DB.GetRoles(ctx, ID)
}

func (r *Role) addRole(ctx context.Context, role entity.Role) error {
	return r.DB.CreateRole(ctx, role)
}

func (r *Role) removeRole(ctx context.Context, ID int64) error {
	return r.DB.DeleteRole(ctx, ID)
}
