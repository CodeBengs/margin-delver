package authrepository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	authconstant "margin-delver/modules/auth/auth_constant"
	authdto "margin-delver/modules/auth/auth_dto"
	authmodel "margin-delver/modules/auth/auth_entity/auth_model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const tokenDuration = 24 * time.Hour

func (repository *AuthRepository) Login(
	ctx context.Context,
	request *authdto.LoginRequest,
) (*authdto.LoginResponse, error) {
	username := strings.TrimSpace(request.Username)
	user, err := repository.findByUsername(ctx, username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, authconstant.ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return nil, authconstant.ErrInvalidCredentials
	}

	return &authdto.LoginResponse{
		Token:     generateToken(),
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(tokenDuration),
		User:      toUserResponse(user),
	}, nil
}

func (repository *AuthRepository) findByUsername(ctx context.Context, username string) (*authmodel.User, error) {
	var user authmodel.User
	err := repository.db.WithContext(ctx).
		Where("username = ? AND flag_active = ?", username, true).
		First(&user).
		Error

	return &user, err
}

func generateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}

	return hex.EncodeToString(bytes)
}

func toUserResponse(user *authmodel.User) authdto.UserResponse {
	return authdto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}
}
