package store

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/vector-10/url-shortner/internal/models"
)

const defaultCacheTTL = 1 * time.Hour

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &RedisCache{client: client}
}

func (r *RedisCache) CacheSlug(record *models.URLRecord) error {
    data, err := json.Marshal(record)
    if err != nil {
        return err
    }

    ttl := defaultCacheTTL
    if record.ExpiresAt != nil {
        remaining := time.Until(*record.ExpiresAt)
        if remaining < ttl {
            ttl = remaining
        }
    }

    ctx := context.Background()
    return r.client.Set(ctx, "slug:"+record.Slug, data, ttl).Err()
}

func (r *RedisCache) GetCachedSlug(slug string) (*models.URLRecord, error) {
    ctx := context.Background()

    data, err := r.client.Get(ctx, "slug:"+slug).Bytes()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var record models.URLRecord
    if err := json.Unmarshal(data, &record); err != nil {
        return nil, err
    }
    return &record, nil
}

func (r *RedisCache) InvalidateSlug(slug string) error {
    ctx := context.Background()
    return r.client.Del(ctx, "slug:"+slug).Err()
}
