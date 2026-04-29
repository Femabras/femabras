// femabras/backend/internal/services/redis.go
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Sentinel errors — use errors.Is() everywhere, never .Error() string matching
var (
	ErrDailyAdLimitReached = errors.New("daily ad limit reached")
	ErrRedisUnavailable    = errors.New("redis unavailable")
)

var rdb *redis.Client

// adRewardScript atomically checks the daily ad cap, increments the ad counter,
// and increments the attempt counter — all in a single round-trip.
// This eliminates the read-then-write race condition of the previous implementation.
//
// KEYS[1] = ads_watched:{userID}:{date}
// KEYS[2] = attempts:{userID}:{date}
// ARGV[1] = daily ad limit
// ARGV[2] = TTL in seconds
var adRewardScript = redis.NewScript(`
local watched = tonumber(redis.call('GET', KEYS[1])) or 0
if watched >= tonumber(ARGV[1]) then
    return redis.error_reply('daily ad limit reached')
end
redis.call('INCR', KEYS[1])
redis.call('EXPIRE', KEYS[1], tonumber(ARGV[2]))
local newAttempts = redis.call('INCR', KEYS[2])
redis.call('EXPIRE', KEYS[2], tonumber(ARGV[2]))
return newAttempts
`)

const (
	challengeEventChannel = "femabras:challenge_events"
	ttl26hSeconds         = int(26 * time.Hour / time.Second)
)

func InitRedis(url string) {
	if url == "" {
		rdb = nil
		return
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		// A bad Redis URL is a fatal misconfiguration — surface it immediately
		panic(fmt.Sprintf("redis: failed to parse URL: %v", err))
	}
	rdb = redis.NewClient(opt)
}

func GetRedisClient() *redis.Client {
	return rdb
}

// ── Attempts ─────────────────────────────────────────────────────────────────

func GetOrCreateAttempts(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 5, nil
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)

	val, err := rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		if setErr := rdb.Set(ctx, key, 5, time.Duration(ttl26hSeconds)*time.Second).Err(); setErr != nil {
			return 0, fmt.Errorf("redis set default attempts: %w", setErr)
		}
		return 5, nil
	}
	if err != nil {
		return 0, fmt.Errorf("redis get attempts: %w", err)
	}
	return val, nil
}

func DecrementAndSave(ctx context.Context, userID, date string) (int, error) {
	if rdb == nil {
		return 4, nil
	}
	key := fmt.Sprintf("attempts:%s:%s", userID, date)
	ttl := time.Duration(ttl26hSeconds) * time.Second

	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis exists: %w", err)
	}
	if exists == 0 {
		if err := rdb.Set(ctx, key, 5, ttl).Err(); err != nil {
			return 0, fmt.Errorf("redis set default: %w", err)
		}
	}

	rem, err := rdb.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis decr: %w", err)
	}
	rdb.Expire(ctx, key, ttl)
	return int(rem), nil
}

func LockOnSuccess(ctx context.Context, userID, date string) {
	if rdb != nil {
		rdb.Set(ctx, fmt.Sprintf("attempts:%s:%s", userID, date), 0,
			time.Duration(ttl26hSeconds)*time.Second)
	}
}

// ── Ad Rewards (atomic Lua) ───────────────────────────────────────────────────

// IncrementAttemptAndAdsWatched atomically enforces the daily ad cap and grants
// one extra attempt via a single Redis Lua script — no race condition possible.
func IncrementAttemptAndAdsWatched(ctx context.Context, userID, date string) error {
	if rdb == nil {
		return nil
	}

	adsKey := fmt.Sprintf("ads_watched:%s:%s", userID, date)
	attemptsKey := fmt.Sprintf("attempts:%s:%s", userID, date)

	// Ensure the attempts key exists before the script tries to INCR it
	if _, err := GetOrCreateAttempts(ctx, userID, date); err != nil {
		return err
	}

	_, err := adRewardScript.Run(
		ctx, rdb,
		[]string{adsKey, attemptsKey},
		3,             // daily ad limit
		ttl26hSeconds, // TTL
	).Int()

	if err != nil {
		if err.Error() == "daily ad limit reached" {
			return ErrDailyAdLimitReached
		}
		return fmt.Errorf("redis ad script: %w", err)
	}

	return nil
}

// ── Pub/Sub (SSE broadcast) ───────────────────────────────────────────────────

// PublishChallengeEvent broadcasts a challenge lifecycle event (e.g. "solved")
// to all connected SSE clients via Redis pub/sub. Works across multiple
// backend instances.
func PublishChallengeEvent(ctx context.Context, event string) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	return rdb.Publish(ctx, challengeEventChannel, event).Err()
}

// SubscribeChallengeEvents opens a Redis pub/sub subscription for challenge
// events. The caller must call Close() on the returned *redis.PubSub.
func SubscribeChallengeEvents(ctx context.Context) *redis.PubSub {
	return rdb.Subscribe(ctx, challengeEventChannel)
}

// ── User Auth Cache ───────────────────────────────────────────────────────────

// GetCachedUserStatus returns true when the user's active+verified status is
// cached, allowing the auth middleware to skip a database round-trip.
func GetCachedUserStatus(ctx context.Context, key string) (bool, error) {
	if rdb == nil {
		return false, ErrRedisUnavailable
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return false, err // redis.Nil = cache miss, other errors = real failure
	}
	return val == "1", nil
}

// SetCachedUserStatus marks a user as verified and active in the cache.
func SetCachedUserStatus(ctx context.Context, key string, ttl time.Duration) {
	if rdb != nil {
		rdb.Set(ctx, key, "1", ttl)
	}
}

// InvalidateUserCache removes the auth cache entry for a user — call on logout
// or whenever a user's verified status changes.
func InvalidateUserCache(ctx context.Context, userID string) {
	if rdb != nil {
		rdb.Del(ctx, fmt.Sprintf("user:active:%s", userID))
	}
}
