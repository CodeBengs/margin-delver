package authhandler

import (
	"errors"

	"margin-delver/base/response"
	authconstant "margin-delver/modules/auth/auth_constant"
	authdto "margin-delver/modules/auth/auth_dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, authconstant.ErrCodeInvalidRequest, "Invalid login request")
		return
	}

	result, err := h.service.Login(c.Request.Context(), &req)
	if errors.Is(err, authconstant.ErrInvalidCredentials) {
		response.Unauthorized(c, authconstant.ErrCodeInvalidCredentials, "Invalid username or password")
		return
	}

	if err != nil {
		h.log.SugarLog().Errorf("failed to login: %v", err)
		response.InternalServerError(c, authconstant.ErrCodeLoginFailed, "Failed to login")
		return
	}

	response.Success(c, "Login success", result)
}
