package handlers

import (
	"jwt-auth-service/internal/service"
	"jwt-auth-service/pkg/jwt"
)

type Handlers struct {
	s            *service.Service
	tokenManager jwt.TokenManager
}

func NewHandlers(s *service.Service, tokenManager jwt.TokenManager) *Handlers {
	return &Handlers{s: s, tokenManager: tokenManager}
}
