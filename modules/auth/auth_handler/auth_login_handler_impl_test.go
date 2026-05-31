package authhandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	moduleLib "margin-delver/lib"
	authconstant "margin-delver/modules/auth/auth_constant"
	authdto "margin-delver/modules/auth/auth_dto"

	"github.com/gin-gonic/gin"
)

type authServiceMock struct {
	loginResponse *authdto.LoginResponse
	loginErr      error
	loginRequest  *authdto.LoginRequest
}

func (m *authServiceMock) Login(ctx context.Context, request *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	m.loginRequest = request
	return m.loginResponse, m.loginErr
}

func TestLoginSuccess(t *testing.T) {
	service := &authServiceMock{
		loginResponse: &authdto.LoginResponse{
			Token:     "token",
			TokenType: "Bearer",
			ExpiresAt: time.Now().Add(time.Hour),
			User: authdto.UserResponse{
				ID:       1,
				Username: "DELVERADMIN1",
				Name:     "Delver Administrator",
			},
		},
	}

	recorder := performLoginRequest(service, `{"username":"DELVERADMIN1","password":"delverAdmin1"}`)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	if service.loginRequest == nil {
		t.Fatalf("expected service login request")
	}

	if service.loginRequest.Username != "DELVERADMIN1" {
		t.Fatalf("expected username delegated to service")
	}
}

func TestLoginInvalidRequest(t *testing.T) {
	service := &authServiceMock{}

	recorder := performLoginRequest(service, `{"username":""}`)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}

	if service.loginRequest != nil {
		t.Fatalf("expected service not called")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	service := &authServiceMock{
		loginErr: authconstant.ErrInvalidCredentials,
	}

	recorder := performLoginRequest(service, `{"username":"DELVERADMIN1","password":"wrong"}`)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
	}
}

func TestLoginInternalServerError(t *testing.T) {
	service := &authServiceMock{
		loginErr: errors.New("db down"),
	}

	recorder := performLoginRequest(service, `{"username":"DELVERADMIN1","password":"delverAdmin1"}`)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
	}
}

func performLoginRequest(service *authServiceMock, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(service, moduleLib.NewBaseLog(&moduleLib.AppConfig{AppEnv: "local"}))
	router := gin.New()
	router.POST("/internal/v1/auth/login", handler.Login)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/internal/v1/auth/login",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	return recorder
}
