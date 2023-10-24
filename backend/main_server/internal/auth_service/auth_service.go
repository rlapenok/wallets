package auth_service

import (
	"github.com/rlapenok/wallets/backend/main_server/internal/storages/tokens_storage"
	"github.com/rlapenok/wallets/backend/main_server/internal/token_service"
)

type AuthService struct {
	storage       tokens_storage.Crud
	token_service token_service.TokenServiceAuth
}

func New() *AuthService {
	storage := tokens_storage.New()
	token_service := token_service.New()
	auth_service := AuthService{storage: storage, token_service: token_service}
	return &auth_service
}
