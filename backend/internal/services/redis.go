// femabras/backend/internal/services/redis.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis(url string) {
	if url == "" {
		rdb = nil // fallback mode (local dev without Redis)
		return
	}
	opt, _ := redis.ParseURL(url)
	rdb = redis.NewClient(opt)
}

func GetOrCreateAttempts(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 5, nil
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)

	val, err := rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		// First time this user plays today → initialize with 5 attempts
		rdb.Set(ctx, key, 5, 26*time.Hour)
		return 5, nil
	}
	if err != nil {
		return 0, fmt.Errorf("redis get error: %w", err)
	}
	return val, nil
}

func DecrementAndSave(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 4, nil // fallback for local testing
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)

	// Ensure the key exists with default 5 before decrementing
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis exists error: %w", err)
	}
	if exists == 0 {
		if err := rdb.Set(ctx, key, 5, 26*time.Hour).Err(); err != nil {
			return 0, fmt.Errorf("redis set default error: %w", err)
		}
	}

	rem, err := rdb.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis decr error: %w", err)
	}
	rdb.Expire(ctx, key, 26*time.Hour)
	return int(rem), nil
}

func LockOnSuccess(ctx context.Context, userID, date string) {
	if rdb != nil {
		rdb.Set(ctx, fmt.Sprintf("attempts:%s:%s", userID, date), 0, 26*time.Hour)
	}
}

func IncrementAttemptAndAdsWatched(ctx context.Context, userID, date string) error {
	if rdb == nil {
		return nil
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)
	return rdb.Incr(ctx, key).Err()
}
