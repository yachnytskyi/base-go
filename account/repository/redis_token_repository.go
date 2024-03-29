package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// redisTokenRepository is data/repository implementation
// of the service layer TokenRepository.
type redisTokenRepository struct {
	Redis *redis.Client
}

// NewTokenRepository is a factory for initializing User Repositories.
func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores a refresh token with an expiry time.
func (repository *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	// We will store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage.
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := repository.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to Redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

// DeleteRefreshToken used to delete old refresh tokens.
// Services my access this to revolve tokens.
func (repository *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := repository.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid.
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

// DeleteUserRefreshTokens looks for all tokens beginning with
// userID and scans to delete them in a non-blocking fashion.
func (repository *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	scanIterator := repository.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failsCount := 0

	for scanIterator.Next(ctx) {
		if err := repository.Redis.Del(ctx, scanIterator.Val()).Err(); err != nil {
			log.Printf("Failed to delete the refresh token: %s\n", scanIterator.Val())
			failsCount++
		}
	}

	// Check the last value.
	if err := scanIterator.Err(); err != nil {
		log.Printf("Failed to delete the refresh token: %s\n", scanIterator.Val())
	}

	if failsCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
