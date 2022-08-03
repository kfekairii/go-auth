package auth_domain

import (
	"goauth/user/user_domain"
	"time"
)

type IAuthService interface {
	Signup(u *user_domain.User) error
}

type ITokenService interface {
	NewPairFromUser(u *user_domain.User, prevAccessToken string) (*TokenPair, error)
	ValidateAccessToken(token string) (uint, error)
	ValidateRefreshToken(refreshToken string) (*RefreshToken, error)
}

type ITokenRepository interface {
	SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(userID string, prevTokenID string) error
}

type IPasswordService interface {
	CreateHash(password string) (string, error)
	ComparePasswordAndHash(password string, hash string) (match bool, err error)
}
