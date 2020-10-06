package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Configurations struct {
	Server *ServerConfig
	Mongo  *MongoConfig
	Redis  *RedisConfig
}

type ServerConfig struct {
	Port          string
	Admin         string
	AccessSecret  string
	AccessExpire  int64
	RefreshExpire int64
	RefreshSecret string
}

type MongoConfig struct {
	Username string
	Password string
	Host     string
	DBName   string
	Query    string
}

type RedisConfig struct {
	Host string
}

func Config() *Configurations {
	// load .env file
	env := os.Getenv("GO_ENV")
	err := godotenv.Load(".env")

	if err != nil && env == "dev" {
		log.Fatal("Error loading .env file")
	}
	// server config
	atExpire, err := strconv.ParseInt(os.Getenv("ACCESS_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	rtExpire, err := strconv.ParseInt(os.Getenv("REFRESH_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	serverCnf := &ServerConfig{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		AccessExpire:  atExpire,
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		RefreshExpire: rtExpire,
		Port:          os.Getenv("SERVER_PORT"),
		Admin:         os.Getenv("ADMIN_MAIL"),
	}
	// mongo config
	mongoCnf := &MongoConfig{
		Username: os.Getenv("MONGO_DB_USERNAME"),
		Password: os.Getenv("MONGO_DB_PASSWORD"),
		Host:     os.Getenv("MONGO_DB_HOST"),
		DBName:   os.Getenv("MONGO_DB_NAME"),
		Query:    os.Getenv("MONGO_DB_QUERY"),
	}

	// redis config

	redisCnf := &RedisConfig{
		Host: os.Getenv("REDIS_DB_HOST"),
	}

	return &Configurations{
		Server: serverCnf,
		Mongo:  mongoCnf,
		Redis:  redisCnf,
	}
}
