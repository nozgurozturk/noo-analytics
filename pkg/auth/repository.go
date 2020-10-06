package auth

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"time"
)

type Repository interface {
	DeleteToken(token *entities.Token) *errors.ApplicationError
	SetToken(token *entities.Token) *errors.ApplicationError
}

type repository struct {
	Keys *redis.Client
}

func NewRepository(keys *redis.Client) Repository {
	return &repository{Keys: keys}
}

func (r *repository) DeleteToken(token *entities.Token) *errors.ApplicationError {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := r.Keys.Del(ctx, token.Uuid).Result()
	if err != nil {
		return errors.InternalServer(err.Error())
	}
	return nil
}

func (r *repository) SetToken(token *entities.Token) *errors.ApplicationError {
	expires := time.Unix(token.Expires, 0)
	now := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := r.Keys.Set(ctx, token.Uuid, token.UserId, expires.Sub(now)).Err()
	if err != nil {
		return errors.InternalServer("Error when setting token")
	}
	return nil
}
