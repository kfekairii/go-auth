package auth_service

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"goauth/utils/apperrors"
	"log"
	"strconv"
	"time"

	"goauth/auth/auth_domain"
	"goauth/user/user_domain"
)

type TokenService struct {
	PvKey           *rsa.PrivateKey
	PubKey          *rsa.PublicKey
	RefreshSecret   string
	TokenRepository auth_domain.ITokenRepository
}

type TSConfig struct {
	PvKey           *rsa.PrivateKey
	PubKey          *rsa.PublicKey
	RefreshSecret   string
	TokenRepository auth_domain.ITokenRepository
}

func NewTokenService(c *TSConfig) auth_domain.ITokenService {
	return &TokenService{
		PvKey:           c.PvKey,
		PubKey:          c.PubKey,
		RefreshSecret:   c.RefreshSecret,
		TokenRepository: c.TokenRepository,
	}
}

func (s *TokenService) NewPairFromUser(u *user_domain.User, prevAccessTokenID string) (*auth_domain.TokenPair, error) {

	if prevAccessTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(strconv.Itoa(int(u.ID)), prevAccessTokenID); err != nil {
			log.Printf("Error Deleting RefreshTokenData ID for user uid: %v, Error:%v\n", u.ID, err)
			return nil, err
		}
	}

	accessToken, err := generateAccessToken(u, s.PvKey)
	if err != nil {
		log.Printf("Error generating AccessToken for user uid: %v, Error:%v\n", u.ID, err)
		return nil, apperrors.NewInternal()
	}
	refreshToken, err := generateRefreshToken(u.ID, s.RefreshSecret)

	if err != nil {
		log.Printf("Error generating RefreshTokenData for user uid: %v, Error:%v\n", u.ID, err)
		return nil, apperrors.NewInternal()
	}
	// Store the Refresh Token
	if err = s.TokenRepository.SetRefreshToken(strconv.Itoa(int(u.ID)), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing RefreshTokenData ID for user uid: %v, Error:%v\n", u.ID, err)
		return nil, apperrors.NewInternal()
	}

	return &auth_domain.TokenPair{
		AccessToken: auth_domain.AccessToken{SS: accessToken},
		RefreshToken: auth_domain.RefreshToken{
			SS:  refreshToken.SS,
			ID:  refreshToken.ID,
			UID: u.ID,
		},
	}, nil
}

func (s *TokenService) ValidateAccessToken(tokenString string) (uint, error) {
	claims, err := validateAccessToken(tokenString, s.PubKey)
	if err != nil {
		return 0, apperrors.NewAuthorization("bad or missing token")
	}
	return claims.UID, nil
}

func validateAccessToken(tokenString string, key *rsa.PublicKey) (*accessTokenClaims, error) {
	claims := &accessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")

	}
	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse access token claims")
	}
	return claims, nil
}

func (s *TokenService) ValidateRefreshToken(refreshToken string) (*auth_domain.RefreshToken, error) {
	claims, err := validateRefreshToken(refreshToken, s.RefreshSecret)
	if err != nil {
		return nil, apperrors.NewAuthorization("bad or missing token")
	}
	return &auth_domain.RefreshToken{
		SS:  refreshToken,
		ID:  claims.ID,
		UID: claims.UID,
	}, nil
}

func validateRefreshToken(tokenString string, key string) (*RefreshTokenClaims, error) {
	claims := &RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")

	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse refresh token claims")
	}

	return claims, nil
}

type accessTokenClaims struct {
	UID uint `json:"uid"`
	jwt.RegisteredClaims
}

func generateAccessToken(u *user_domain.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now()
	expiredAt := unixTime.Add(time.Minute * 15) // expired after 15 minutes

	claims := accessTokenClaims{
		UID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to Sign AccessToken")
		return "", err
	}

	return ss, nil
}

type RefreshTokenData struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

type RefreshTokenClaims struct {
	UID uint `json:"uid"`
	jwt.RegisteredClaims
}

func generateRefreshToken(uid uint, key string) (*RefreshTokenData, error) {
	currentTime := time.Now()
	expiredAt := currentTime.AddDate(0, 0, 7)
	tokenId, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token random id")
		return nil, err
	}

	claims := &RefreshTokenClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			ID:        tokenId.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token")
		return nil, err
	}

	return &RefreshTokenData{
		SS:        ss,
		ID:        tokenId.String(),
		ExpiresIn: expiredAt.Sub(currentTime),
	}, nil
}
