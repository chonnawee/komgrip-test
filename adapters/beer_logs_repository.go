package adapters

import (
	"context"
	"komgrip-test/entities"
	"komgrip-test/usecases"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type beerLogsRepositoryDB struct {
	collection *mongo.Collection
}

func NewBeerLogsRepositoryDB(collection *mongo.Collection) usecases.BeerLogsRepository {
	return &beerLogsRepositoryDB{collection: collection}
}

func (r *beerLogsRepositoryDB) CreateLog(log entities.BeerLogs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	doc := bson.M{
		"_id":          primitive.NewObjectID(),
		"method":       log.Method,
		"request_data": log.RequestData,
		"status":       log.Status,
		"created_at":   log.CreatedAt,
	}
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
