package blueprints

var JWTManager = `package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
	minimumSecretSize = 32
)

type Config struct {
	Issuer        string
	Audience      string
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Claims struct {
	Roles     []string ` + "`json:\"roles,omitempty\"`" + `
	TokenType string   ` + "`json:\"token_type\"`" + `
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	TokenType    string ` + "`json:\"token_type\"`" + `
	ExpiresIn    int64  ` + "`json:\"expires_in\"`" + `
}

type Manager struct{ config Config }

func NewManager(config Config) (*Manager, error) {
	if len(config.AccessSecret) < minimumSecretSize || len(config.RefreshSecret) < minimumSecretSize {
		return nil, fmt.Errorf("JWT access and refresh secrets must each contain at least %d bytes", minimumSecretSize)
	}
	if config.AccessSecret == config.RefreshSecret {
		return nil, errors.New("JWT access and refresh secrets must be different")
	}
	if config.Issuer == "" || config.Audience == "" {
		return nil, errors.New("JWT issuer and audience are required")
	}
	if config.AccessTTL <= 0 || config.RefreshTTL <= 0 || config.AccessTTL >= config.RefreshTTL {
		return nil, errors.New("JWT TTL values must be positive and access TTL must be shorter than refresh TTL")
	}
	return &Manager{config: config}, nil
}

func NewManagerFromViper() (*Manager, error) {
	return NewManager(Config{
		Issuer: viper.GetString("JWT_ISSUER"), Audience: viper.GetString("JWT_AUDIENCE"),
		AccessSecret: viper.GetString("JWT_ACCESS_SECRET"), RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
		AccessTTL: viper.GetDuration("JWT_ACCESS_TTL"), RefreshTTL: viper.GetDuration("JWT_REFRESH_TTL"),
	})
}

func (m *Manager) IssuePair(subject string, roles []string) (TokenPair, error) {
	if subject == "" { return TokenPair{}, errors.New("JWT subject is required") }
	now := time.Now().UTC()
	access, err := m.sign(subject, roles, AccessTokenType, now, m.config.AccessTTL, m.config.AccessSecret)
	if err != nil { return TokenPair{}, err }
	refresh, err := m.sign(subject, roles, RefreshTokenType, now, m.config.RefreshTTL, m.config.RefreshSecret)
	if err != nil { return TokenPair{}, err }
	return TokenPair{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: int64(m.config.AccessTTL.Seconds())}, nil
}

func (m *Manager) Refresh(refreshToken string) (TokenPair, error) {
	claims, err := m.ParseRefresh(refreshToken)
	if err != nil { return TokenPair{}, err }
	return m.IssuePair(claims.Subject, claims.Roles)
}

func (m *Manager) ParseAccess(token string) (*Claims, error) {
	return m.parse(token, AccessTokenType, m.config.AccessSecret)
}

func (m *Manager) ParseRefresh(token string) (*Claims, error) {
	return m.parse(token, RefreshTokenType, m.config.RefreshSecret)
}

func (m *Manager) sign(subject string, roles []string, tokenType string, now time.Time, ttl time.Duration, secret string) (string, error) {
	id, err := tokenID()
	if err != nil { return "", fmt.Errorf("create token ID: %w", err) }
	claims := Claims{Roles: roles, TokenType: tokenType, RegisteredClaims: jwt.RegisteredClaims{
		Issuer: m.config.Issuer, Subject: subject, Audience: jwt.ClaimStrings{m.config.Audience}, ID: id,
		IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func (m *Manager) parse(raw, expectedType, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 { return nil, fmt.Errorf("unexpected JWT signing method %q", token.Method.Alg()) }
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(m.config.Issuer),
		jwt.WithAudience(m.config.Audience), jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil { return nil, fmt.Errorf("invalid %s token: %w", expectedType, err) }
	if !token.Valid { return nil, fmt.Errorf("invalid %s token", expectedType) }
	if claims.TokenType != expectedType { return nil, fmt.Errorf("expected %s token", expectedType) }
	return claims, nil
}

func tokenID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil { return "", err }
	return hex.EncodeToString(value), nil
}
`

var JWTMiddleware = `package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"{{ .ModuleName }}/pkg/auth"
)

const JWTClaimsLocal = "jwt_claims"

func RequireJWT(manager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid bearer token")
		}
		claims, err := manager.ParseAccess(parts[1])
		if err != nil { return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired access token") }
		c.Locals(JWTClaimsLocal, claims)
		return c.Next()
	}
}

func JWTClaims(c *fiber.Ctx) (*auth.Claims, bool) {
	claims, ok := c.Locals(JWTClaimsLocal).(*auth.Claims)
	return claims, ok
}
`
