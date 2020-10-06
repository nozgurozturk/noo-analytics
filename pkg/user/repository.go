package user

import (
	"context"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Create(u *entities.User) (user *entities.User, error *errors.ApplicationError)
	FindByEmail(email string) (user *entities.User, error *errors.ApplicationError)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) FindByEmail(email string) (user *entities.User, error *errors.ApplicationError) {
	u := new(entities.User)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)

	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	return u, nil
}

func (r *repository) Create(u *entities.User) (user *entities.User, error *errors.ApplicationError) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := r.Collection.InsertOne(ctx, &u)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	return u, nil
}
