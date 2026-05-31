package authhandler

import (
	moduleLib "margin-delver/lib"
	authservice "margin-delver/modules/auth/auth_service"
)

type Handler struct {
	service authservice.AuthServiceInterface
	log     *moduleLib.BaseLog
}

func NewHandler(service authservice.AuthServiceInterface, log *moduleLib.BaseLog) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}
