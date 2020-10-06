package analytics

import (
	"context"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	FindActionsByDate(entity *entities.AnalyticsActionRequest) ([]entities.AnalyticsActionResponse, *errors.ApplicationError)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) FindActionsByDate(entity *entities.AnalyticsActionRequest) ([]entities.AnalyticsActionResponse, *errors.ApplicationError) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var matchFilter bson.D
	var groupByFilter bson.D
	var groupByObject bson.D

	matchFilter = append(matchFilter, bson.E{Key: "act", Value: entity.Action})

	if entity.Year.IsInclude != false {
		groupByFilter = append(groupByFilter, bson.E{Key: "y", Value: "$y"})
		groupByObject = append(groupByObject, bson.E{Key: "y", Value: bson.M{"$first": "$y"}})

		to := int(time.Now().Year()) + 1
		if entity.Year.Range.To != 0 {
			to = entity.Year.Range.To
		}
		matchFilter = append(matchFilter, bson.E{Key: "y", Value: bson.M{
			"$gt": entity.Year.Range.From - 1,
			"$lt": to,
		}})
	}

	if entity.Month.IsInclude != false {
		groupByFilter = append(groupByFilter, bson.E{Key: "m", Value: "$m"})
		groupByObject = append(groupByObject, bson.E{Key: "m", Value: bson.M{"$first": "$m"}})

		to := int(time.Now().Year()) + 1
		if entity.Month.Range.To != 0 {
			to = entity.Month.Range.To
		}
		matchFilter = append(matchFilter, bson.E{Key: "m", Value: bson.M{
			"$gt": entity.Year.Range.From - 1,
			"$lt": to,
		}})
	}

	if entity.Day.IsInclude != false {
		groupByFilter = append(groupByFilter, bson.E{Key: "d", Value: "$d"})
		groupByObject = append(groupByObject, bson.E{Key: "d", Value: bson.M{"$first": "$d"}})

		to := int(time.Now().Year()) + 1
		if entity.Day.Range.To != 0 {
			to = entity.Day.Range.To
		}
		matchFilter = append(matchFilter, bson.E{Key: "d", Value: bson.M{
			"$gt": entity.Day.Range.From - 1,
			"$lt": to,
		}})
	}

	if entity.Hour.IsInclude != false {
		groupByFilter = append(groupByFilter, bson.E{Key: "h", Value: "$h"})
		groupByObject = append(groupByObject, bson.E{Key: "h", Value: bson.M{"$first": "$h"}})

		to := int(time.Now().Year()) + 1
		if entity.Hour.Range.To != 0 {
			to = entity.Hour.Range.To
		}
		matchFilter = append(matchFilter, bson.E{Key: "d", Value: bson.M{
			"$gt": entity.Hour.Range.From - 1,
			"$lt": to,
		}})
	}

	matchStage := bson.D{{
		"$match",
		matchFilter,
	}}
	groupByObject = append(groupByObject, bson.E{Key: "v", Value: bson.D{{"$push", "$ip"}}})
	groupByObject = append(groupByObject, bson.E{Key: "uv", Value: bson.D{{"$addToSet", "$ip"}}})
	groupByObject = append(groupByObject, bson.E{Key: "_id", Value: groupByFilter})

	groupStage := bson.D{{
		"$group", groupByObject,
	},
	}

	sortState := bson.D{{
		"$sort", bson.M{
			"ts": 1,
		},
	}}

	findActionCursor, err := r.Collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortState})

	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	var analytics []entities.AnalyticsActionResponse
	if err = findActionCursor.All(ctx, &analytics); err != nil {
		panic(err)
	}

	return analytics, nil
}
