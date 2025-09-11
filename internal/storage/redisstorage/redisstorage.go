package redisstorage

import (
	"Gonoty/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *redis.Client
	queue  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisStorage(conf RedisConfig) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return nil, err
	}

	fmt.Println("Connected to Redis:", pong)
	return &RedisStorage{
		Client: client,
		queue:  "email_queue"}, nil
}

func (r *RedisStorage) Add(ctx context.Context, task models.Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	err = r.Client.LPush(ctx, r.queue, taskJSON).Err()
	if err != nil {
		return fmt.Errorf("redis Lpush error: %v", err)
	}

	return nil
}

func (r *RedisStorage) List(ctx context.Context) error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("List ping error: %v", err)
	}

	length, err := r.Client.LLen(ctx, r.queue).Result()
	if err != nil {
		return fmt.Errorf("List len error: %v", err)
	}

	fmt.Printf("Queue length: %d\n", length)

	if length > 0 {
		items, err := r.Client.LRange(ctx, r.queue, 0, 4).Result()
		if err != nil {
			log.Fatal("Error reading queue:", err)
		}

		fmt.Println("\nFirst 5 items in queue:")
		for i, item := range items {
			var task map[string]interface{}
			if err := json.Unmarshal([]byte(item), &task); err == nil {
				fmt.Printf("%d: %+v\n", i+1, task)
			} else {
				fmt.Printf("%d: %s\n", i+1, item)
			}
		}
	}
	return nil
}

func (r *RedisStorage) Delete(ctx context.Context, taskID string) error {
	// deleted, err := r.client. (ctx, r.hash, taskID).Result()
	// if err != nil {
	// 	return err
	// }

	// if deleted == 0 {
	// 	return fmt.Errorf("task not found: %s", taskID)
	// }

	return nil
}
