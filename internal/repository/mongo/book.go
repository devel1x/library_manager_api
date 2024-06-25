package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/validator"
	"time"
)

func (r *MongoRepo) GetBookByISBN(ctx context.Context, id string) (*entity.Book, error) {
	filter := bson.D{{"_id", id}}
	var book entity.Book
	err := r.booksCollection.FindOne(ctx, filter).Decode(&book)
	switch {
	case err == nil:
		return &book, nil
	case err == mongo.ErrNoDocuments:
		return nil, utils.ErrNotExist
	default:
		return nil, err
	}
	return &book, nil
}

func (r *MongoRepo) CreateBook(ctx context.Context, book *entity.BookForm) (interface{}, error) {
	res, err := r.booksCollection.InsertOne(ctx, book)
	if err != nil {
		fmt.Println(err)
		if mongo.IsDuplicateKeyError(err) {
			return nil, utils.ErrBookAlreadyExists
		}
		return nil, err
	}
	return res.InsertedID, nil
}
func (r *MongoRepo) ListBook(ctx context.Context, page, pageSize int) (*entity.PaginatedBooks, error) {
	var result entity.PaginatedBooks

	totalCount, err := r.booksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %v", err)
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if lastPage == 0 {
		return &entity.PaginatedBooks{nil, 0}, nil
	}
	if page > lastPage {
		return nil, fmt.Errorf("the last page is %d: %w", lastPage, utils.ErrBadInput)
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)

	cursor, err := r.booksCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch documents: %v", err)
	}
	defer cursor.Close(ctx)

	var books []*entity.Book
	for cursor.Next(ctx) {
		var book entity.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		books = append(books, &book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	result.Books = books
	result.LastPage = lastPage
	return &result, nil
}

func (r *MongoRepo) UpdateBookByISBN(ctx context.Context, book *entity.BookForm) error {
	updateQuery := bson.M{}

	updateQuery["updatedAt"] = time.Now()

	if validator.MinChars(book.Title, 1) {
		updateQuery["title"] = book.Title
	}

	if validator.MinChars(book.Publisher, 1) {
		updateQuery["publisher"] = book.Publisher
	}

	if validator.CheckArr(book.Author) {
		updateQuery["author"] = book.Author
	}

	result, err := r.booksCollection.UpdateOne(ctx,
		bson.M{"_id": book.ISBN}, bson.M{"$set": updateQuery})

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return utils.ErrNotExist
	}
	return nil
}
func (r *MongoRepo) DeleteBookByISBN(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	result, err := r.booksCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return utils.ErrNotExist
	}
	return nil
}
