package authservice

import (
	"context"
	"errors"
	"testing"
	"time"

	moduleLib "margin-delver/lib"
	authconstant "margin-delver/modules/auth/auth_constant"
	authdto "margin-delver/modules/auth/auth_dto"
	authmodel "margin-delver/modules/auth/auth_entity/auth_model"
)

type authRepositoryMock struct {
	loginResponse *authdto.LoginResponse
	loginErr      error
	loginRequest  *authdto.LoginRequest
}

func (m *authRepositoryMock) Login(ctx context.Context, request *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	m.loginRequest = request
	return m.loginResponse, m.loginErr
}

func (m *authRepositoryMock) Create(ctx context.Context, user *authmodel.User) error {
	return nil
}

func (m *authRepositoryMock) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	expected := &authdto.LoginResponse{
		Token:     "token",
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(time.Hour),
		User: authdto.UserResponse{
			ID:       1,
			Username: "DELVERADMIN1",
			Name:     "Delver Administrator",
		},
	}
	repository := &authRepositoryMock{
		loginResponse: expected,
	}
	service := NewService(repository, moduleLib.NewBaseLog(&moduleLib.AppConfig{AppEnv: "local"}))
	request := &authdto.LoginRequest{
		Username: "DELVERADMIN1",
		Password: "delverAdmin1",
	}

	actual, err := service.Login(context.Background(), request)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if actual != expected {
		t.Fatalf("expected same response pointer")
	}

	if repository.loginRequest != request {
		t.Fatalf("expected request delegated to repository")
	}
}

func TestAuthServiceLoginError(t *testing.T) {
	expectedErr := authconstant.ErrInvalidCredentials
	repository := &authRepositoryMock{
		loginErr: expectedErr,
	}
	service := NewService(repository, moduleLib.NewBaseLog(&moduleLib.AppConfig{AppEnv: "local"}))

	actual, err := service.Login(context.Background(), &authdto.LoginRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if actual != nil {
		t.Fatalf("expected nil response, got %#v", actual)
	}
}
