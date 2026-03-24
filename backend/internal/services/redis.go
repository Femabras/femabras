// femabras/backend/internal/services/redis.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	// "gorm.io/gorm"
)

var rdb *redis.Client

func InitRedis(url string) {
	if url == "" {
		rdb = nil // fallback mode
		return
	}
	opt, _ := redis.ParseURL(url)
	rdb = redis.NewClient(opt)
}

func GetOrCreateAttempts(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 5, nil // DB fallback
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)
	val, err := rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		rdb.Set(ctx, key, 5, 26*time.Hour)
		return 5, nil
	}
	return val, err
}

func DecrementAndSave(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 4, nil // fake fallback
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)
	rem, _ := rdb.Decr(ctx, key).Result()
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
		return nil // Fake fallback for local testing when REDIS_URL is empty
	}

	key := fmt.Sprintf("attempts:%s:%s", userID, date)

	// Increment the attempts by 1
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	// You could also track 'ads_watched' here in another Redis key if you want a daily limit!
	return nil
}
