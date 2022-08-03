package auth_repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goauth/auth/auth_domain"
	"goauth/utils/apperrors"
	"log"
	"time"
)

type RedisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) auth_domain.ITokenRepository {
	return &RedisTokenRepository{Redis: redisClient}
}

func (r RedisTokenRepository) SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	err := r.Redis.Set(context.Background(), key, 101, expiresIn).Err()
	if err != nil {

		return apperrors.NewInternal()
	}

	return nil
}

func (r RedisTokenRepository) DeleteRefreshToken(userID string, prevTokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, prevTokenID)
	result := r.Redis.Del(context.Background(), key)
	if err := result.Err(); err != nil {
		return apperrors.NewInternal()
	}
	log.Printf("VALLllllll = %v", result.Val())
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID: %s\n doesn't exists\n", userID)
		return apperrors.NewAuthorization("Invalid refresh token\n")
	}
	return nil
}
