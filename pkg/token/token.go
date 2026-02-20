package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"mini_jira/config"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AccessClaims struct {
	UserID    uint   `json:"uid"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID    uint   `json:"uid"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

func NewManagerFromConfig(cfg *config.AppConfig) (*Manager, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	privateKey, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, err
	}
	publicKey, err := loadPublicKey(cfg)
	if err != nil {
		return nil, err
	}

	return &Manager{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     cfg.JWTIssuer,
		audience:   cfg.JWTAudience,
		accessTTL:  time.Duration(cfg.JWTAccessTTLMinutes) * time.Minute,
		refreshTTL: time.Duration(cfg.JWTRefreshTTLDays) * 24 * time.Hour,
	}, nil
}

func (m *Manager) NewAccessToken(userID uint, username, role string) (string, error) {
	now := time.Now().UTC()
	claims := AccessClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Audience:  []string{m.audience},
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(m.privateKey)
}

func (m *Manager) NewRefreshToken(userID uint) (string, error) {
	now := time.Now().UTC()
	claims := RefreshClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Audience:  []string{m.audience},
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(m.privateKey)
}

func (m *Manager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *Manager) RefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *Manager) VerifyAccessToken(tokenString string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, m.keyFunc, jwt.WithAudience(m.audience), jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func (m *Manager) VerifyRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, m.keyFunc, jwt.WithAudience(m.audience), jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func (m *Manager) keyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return m.publicKey, nil
}

func loadPrivateKey(cfg *config.AppConfig) (*rsa.PrivateKey, error) {
	var data []byte
	if cfg.JWTPrivateKeyPath != "" {
		loaded, err := os.ReadFile(cfg.JWTPrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key file: %w", err)
		}
		data = loaded
	} else if cfg.JWTPrivateKey != "" {
		data = []byte(cfg.JWTPrivateKey)
	} else {
		return nil, errors.New("private key is not configured")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	return privateKey, nil
}

func loadPublicKey(cfg *config.AppConfig) (*rsa.PublicKey, error) {
	var data []byte
	if cfg.JWTPublicKeyPath != "" {
		loaded, err := os.ReadFile(cfg.JWTPublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read public key file: %w", err)
		}
		data = loaded
	} else if cfg.JWTPublicKey != "" {
		data = []byte(cfg.JWTPublicKey)
	} else {
		return nil, errors.New("public key is not configured")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	return publicKey, nil
}
