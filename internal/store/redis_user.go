package store

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/redis/go-redis/v9"
    "github.com/vector-10/url-shortner/internal/models"
)

type RedisUserStore struct {
	client *redis.Client
}

func NewRedisUserStore(addr string) *RedisUserStore {
	client := redis.NewClient(&redis.Options{Addr: addr})
    return &RedisUserStore{client: client}
}

func (r *RedisUserStore) CreateUser(user *models.User) error {
	ctx := context.Background()

	exists, err := r.client.Exists(ctx, "user:email:"+user.Email).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("email already registered")
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, "user:"+user.ID, data, 0).Err(); err != nil {
		return err
	}
	return r.client.Set(ctx, "user:email:"+user.Email, user.ID, 0).Err()
}

func (r *RedisUserStore) GetUserByEmail(email string) (*models.User, error) {
    ctx := context.Background()

    id, err := r.client.Get(ctx, "user:email:"+email).Result()
    if err == redis.Nil {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }

    return r.GetUserByID(id)
}

func (r *RedisUserStore) GetUserByID(id string) (*models.User, error) {
    ctx := context.Background()

    data, err := r.client.Get(ctx, "user:"+id).Bytes()
    if err == redis.Nil {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }

    var user models.User
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, err
    }
    return &user, nil
}