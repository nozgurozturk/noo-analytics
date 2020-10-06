package trace

import (
	"context"
	"github.com/nozgurozturk/noo-analytics/entities"
	"github.com/nozgurozturk/noo-analytics/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Insert(entity *entities.Trace) (*entities.Trace, *errors.ApplicationError)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) Insert(entity *entities.Trace) (*entities.Trace, *errors.ApplicationError) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := r.Collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	return entity, nil
}
