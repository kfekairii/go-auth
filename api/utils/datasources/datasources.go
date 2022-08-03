package datasources

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goauth/user/user_domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func Init() (*DataSources, error) {
	// Postgres
	dbURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't connect to database", err)

	}

	err = db.AutoMigrate(&user_domain.User{})

	if err != nil {
		log.Fatalln("Couldn't migrate", err)
	}

	// Redis
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	log.Printf("Connecting to redis\n")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:       0,
		Password: "",
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}
	log.Printf("Connected to redis\n")

	return &DataSources{
		DB:          db,
		RedisClient: rdb,
	}, nil
}

func (d *DataSources) Close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing redis connextion: %w", err)
	}
	return nil
}
