package authrepository

import (
	"context"

	authmodel "margin-delver/modules/auth/auth_entity/auth_model"
)

func (r *AuthRepository) Create(ctx context.Context, user *authmodel.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *AuthRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&authmodel.User{}).
		Count(&total).
		Error

	return total, err
}
