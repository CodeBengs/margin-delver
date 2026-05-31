package authservice

import (
	"context"

	authdto "margin-delver/modules/auth/auth_dto"
)

func (service *AuthService) Login(
	ctx context.Context,
	request *authdto.LoginRequest,
) (*authdto.LoginResponse, error) {
	response, err := service.authRepository.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
