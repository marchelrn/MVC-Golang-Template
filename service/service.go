package service

import (
	"fmt"
	"mini_jira/config"
	"mini_jira/contract"
	"mini_jira/pkg/token"
)

func New(repo *contract.Repository) (*contract.Service, error) {
	cfg := config.GetConfig()
	tokenManager, err := token.NewManagerFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("init token manager: %w", err)
	}

	return &contract.Service{
		User:              ImplUserService(repo.User, repo.RefreshToken, repo.EmailVerification, tokenManager, cfg.JWTRefreshHashSecret),
		EmailVerification: ImplEmailVerificationService(repo.EmailVerification, repo.User),
		Token:             tokenManager,
	}, nil
}
