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

	return r.client.Set(ctx, record.Slug, data, ttl).Err()
}

func (r *RedisStore) GetBySlug(slug string) (*models.URLRecord, error) {
	ctx := context.Background()
	data, err := r.client.Get(ctx, slug).Bytes()
	if err == redis.Nil {
		return nil, errors.New("Key not found")
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

func (r *RedisStore) IncrementClicks(slug string) error {
	record, err := r.GetBySlug(slug)
	if err != nil {
		return err
	}

	record.Clicks++
	return r.Save(record)
}


func (r *RedisStore) Delete(slug string) error {
	ctx := context.Background()
	return r.client.Del(ctx, slug).Err()
}