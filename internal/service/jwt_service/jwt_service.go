package service

import (
	"errors"
	model "tic-tac-toe/internal/domain/model/user"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL   = 15 * time.Minute   // короткий срок жизни
	RefreshTokenTTL  = 7 * 24 * time.Hour // 7 дней
	Issuer           = "tic-tac-toe"      // кто выдал токен
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtProvider struct {
	secretKey []byte
}

func NewJwtProvider(secretKey []byte) JwtProvider {
	return &jwtProvider{
		secretKey: secretKey,
	}
}
func (p *jwtProvider) generateToken(userID uuid.UUID, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Subject:   tokenType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}
func (p *jwtProvider) GenerateAccessToken(user model.User) (string, error) {
	return p.generateToken(user.UUID, TokenTypeAccess, AccessTokenTTL)
}

func (p *jwtProvider) GenerateRefreshToken(user model.User) (string, error) {
	return p.generateToken(user.UUID, TokenTypeRefresh, RefreshTokenTTL)
}

func (p *jwtProvider) parseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil
	})
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

func (p *jwtProvider) ValidateAccessToken(tokenStr string) (*CustomClaims, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if claims.Subject != TokenTypeAccess {
		return nil, errors.New("token is not a access token")
	}
	return claims, nil
}

func (p *jwtProvider) ValidateRefreshToken(tokenStr string) (*CustomClaims, error) {
	claims, err := p.parseToken(tokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if claims.Subject != TokenTypeRefresh {
		return nil, errors.New("token is not a refresh token")
	}
	return claims, nil
}

func (p *jwtProvider) GetUUIDFromToken(token string) (uuid.UUID, error) {
	claims, err := p.parseToken(token)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid token")
	}

	return claims.UserID, nil
}
