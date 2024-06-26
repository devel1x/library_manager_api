package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/validator"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisRepo(client *redis.Client, ttl time.Duration) *RedisRepo {
	return &RedisRepo{
		client: client,
		ttl:    ttl,
	}
}

func bookIDKey(id string) string {
	return fmt.Sprintf("book:%s", id)
}

func (r *RedisRepo) InsertBook(ctx context.Context, book *entity.Book) error {
	data, err := json.Marshal(book)
	fmt.Println("book:", book)
	if err != nil {
		return fmt.Errorf("failed to encode book: %w", err)
	}

	key := bookIDKey(book.ISBN)

	txn := r.client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), r.ttl)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "books", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to books set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) FindBookByISBN(ctx context.Context, id string) (*entity.Book, error) {
	key := bookIDKey(id)

	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, utils.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("get book: %w", err)
	}

	var book entity.Book
	err = json.Unmarshal([]byte(value), &book)
	if err != nil {
		return nil, fmt.Errorf("failed to decode book json: %w", err)
	}

	return &book, nil
}

func (r *RedisRepo) DeleteBookByISBN(ctx context.Context, id string) error {
	key := bookIDKey(id)

	txn := r.client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return utils.ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get book: %w", err)
	}

	if err := txn.SRem(ctx, "books", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from books set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) UpdateBookByISBN(ctx context.Context, book *entity.BookFormUpdate) error {
	key := bookIDKey(book.ISBN)
	existingData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return utils.ErrNotExist
		}
		return fmt.Errorf("get book: %w", err)
	}

	var existingBook map[string]interface{}
	if err := json.Unmarshal([]byte(existingData), &existingBook); err != nil {
		return fmt.Errorf("unmarshal existing book: %w", err)
	}

	if validator.MinChars(book.Title, 1) {
		existingBook["title"] = book.Title
	}
	if validator.MinChars(book.Publisher, 1) {
		existingBook["publisher"] = book.Publisher
	}
	if validator.CheckArr(book.Author) {
		existingBook["author"] = book.Author
	}
	existingBook["updatedAt"] = time.Now().Format(time.RFC3339)

	// Marshal the updated book data back to JSON
	updatedData, err := json.Marshal(existingBook)
	if err != nil {
		return fmt.Errorf("marshal updated book: %w", err)
	}

	// Set the updated book data in Redis
	err = r.client.Set(ctx, key, updatedData, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("set book: %w", err)
	}

	return nil
}
