package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Database struct {
	mongo *mongo.Database
	redis *redis.Client
}

func MongoConnect(cnf *config.MongoConfig) (*mongo.Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s", cnf.Username, cnf.Password, cnf.Host, cnf.DBName, cnf.Query)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Mongo is successfully connected to Application")

	db := client.Database(cnf.DBName)

	return db, err
}

func RedisConnect(cnf *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cnf.Address,
		Username: cnf.UserName,
		Password: cnf.Password,
		DB:       cnf.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Redis is successfully connected to Application")
	return client, nil
}
