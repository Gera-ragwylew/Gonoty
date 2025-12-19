package redisstorage

import (
	"Gonoty/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	queue  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var config = RedisConfig{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func NewRedisStorage() *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisStorage{
		client: client,
		queue:  "email_queue"}
}

func (r *RedisStorage) CheckStatus(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	pong, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis:", pong)
	return nil
}

func (r *RedisStorage) Enqueue(ctx context.Context, task models.Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	err = r.client.LPush(ctx, r.queue, taskJSON).Err()
	if err != nil {
		return fmt.Errorf("redis Lpush error: %v", err)
	}

	return nil
}

func (r *RedisStorage) Dequeue(ctx context.Context) (models.Task, error) {
	result, err := r.client.BRPop(ctx, time.Second, "email_queue").Result()
	if err != nil {
		if err == redis.Nil {
			return models.Task{}, fmt.Errorf("queue empty: timeout")
		} else {
			return models.Task{}, fmt.Errorf("lost connection with redis")
		}
		// log.Printf("Redis error: %v", err)
		// return models.Task{}, err
		// time.Sleep(5 * time.Second)
	}

	var task models.Task
	json.Unmarshal([]byte(result[1]), &task)

	return task, nil
}

func (r *RedisStorage) DequeueBatch(ctx context.Context, batchSize int) ([]models.Task, error) {
	pipe := r.client.Pipeline()

	lrange := pipe.LRange(ctx, "email_queue", 0, int64(batchSize-1))

	ltrim := pipe.LTrim(ctx, "email_queue", int64(batchSize), -1)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	items, err := lrange.Result()
	if err != nil {
		return nil, err
	}

	_, _ = ltrim.Result()

	var tasks []models.Task
	for _, item := range items {
		var task models.Task
		if err := json.Unmarshal([]byte(item), &task); err == nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (r *RedisStorage) List(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("List ping error: %v", err)
	}

	length, err := r.client.LLen(ctx, r.queue).Result()
	if err != nil {
		return fmt.Errorf("List len error: %v", err)
	}

	fmt.Printf("Queue length: %d\n", length)

	if length > 0 {
		items, err := r.client.LRange(ctx, r.queue, 0, 4).Result()
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

func (r *RedisStorage) Close() {
}
