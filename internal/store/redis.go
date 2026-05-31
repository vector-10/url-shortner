package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/vector-10/url-shortner/internal/models"
)


type RedisStore struct {
	client *redis.Client
}


func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{client: client}
}

func (r *RedisStore) Save(record *models.URLRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	ttl := time.Until(record.ExpiresAt)
	ctx := context.Background()

	ok, err := r.client.SetNX(ctx, "slug:"+record.Slug, data, ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("slug already taken")
	}

	 if record.UserID != "" {
        r.client.SAdd(ctx, "user:"+record.UserID+":urls", record.Slug)
    }

	return nil
}

func (r *RedisStore) GetBySlug(slug string) (*models.URLRecord, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, "slug:"+slug).Bytes()
	if err == redis.Nil {
		return nil, errors.New("slug not found")
	}
	if err != nil {
		return nil, err
	}

	var record models.URLRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	clicks, err := r.client.Get(ctx, "clicks:"+slug).Int()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	record.Clicks = clicks

	return &record, nil
}

func (r *RedisStore) ListByUser(userID string) ([]*models.URLRecord, error) {
	ctx := context.Background()

	slugs, err := r.client.SMembers(ctx, "user:"+userID+":urls").Result()
	if err != nil {
		return nil, err
	}

	var records []*models.URLRecord
	for _, slug := range slugs {
		record, err := r.GetBySlug(slug)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *RedisStore) IncrementClicks(slug string) error {
	ctx := context.Background()
	return r.client.Incr(ctx, "clicks:"+slug).Err()
}


func (r *RedisStore) Delete(slug string) error {
	ctx := context.Background()
	return r.client.Del(ctx, "slug:"+slug).Err()
}
