package database

import (
	"context"
	"log"
	"time"

	"hshelby-tkcled-product/config"
	"hshelby-tkcled-product/mongo"
	"hshelby-tkcled-product/src/database/collection"

	src_const "hshelby-tkcled-product/src/const"
)

func ConnectDatabse(ctx context.Context) error {
	var mongoClient *mongo.MongoDB
	var err error
	numberRetry := config.Get().NumberRetry
	if numberRetry == 0 {
		numberRetry = src_const.DEFAULTNUMBERRETRY
	}

	for i := 1; i <= config.Get().NumberRetry; i++ {
		mongoClient, err = mongo.NewMongoDBFromUrl(ctx, config.Get().MongoURL, time.Second*10)
		if err != nil {
			if i == config.Get().NumberRetry {
				log.Println(err)
				return err
			}
			time.Sleep(10 * time.Second)
		}

		if mongoClient != nil {
			break
		}
	}

	if err := collection.LoadProductCollectionMongo(mongoClient); err != nil {
		return err
	}

	return nil
}
