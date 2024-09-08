package collection

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hshelby-tkcled-product/config"
	"hshelby-tkcled-product/src/utilities"

	"hshelby-tkcled-product/mongo"
)

const ProductCollection = "product"

const (
	PrductIndexCategory = "category_id"
)

var (
	_productCollection        *ProductMongoCollection
	loadProductRepositoryOnce sync.Once
)

type ProductMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadProductCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadProductRepositoryOnce.Do(func() {
		_productCollection, err = NewProductMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Product() *ProductMongoCollection {
	if _productCollection == nil {
		panic("database: like product collection is not initiated")
	}
	return _productCollection
}

func NewProductMongoCollection(client *mongo.MongoDB, databaseName string) (*ProductMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewProductMongoCollection] client nil pointer")
	}
	repo := &ProductMongoCollection{
		client:         client,
		collectionName: ProductCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	repo.SetIndex()
	return repo, nil
}

func (repo *ProductMongoCollection) SetIndex() {
	col := repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)

	indexes := []mongoDriver.IndexModel{
		{
			Keys: bson.M{
				PrductIndexCategory: 1,
			},
			Options: &options.IndexOptions{
				Name:   utilities.SetString(PrductIndexCategory),
				Unique: utilities.SetBool(true),
			},
		},
	}

	if !repo.needIndex(col) {
		return
	}

	col.Indexes().CreateMany(context.Background(), indexes)
}

func (repo *ProductMongoCollection) needIndex(col *mongoDriver.Collection) bool {
	keyIndexes := []string{
		PrductIndexCategory,
	}

	listIndexes, err := col.Indexes().ListSpecifications(context.Background())
	if err != nil {
		return true
	}
	indexed := make([]string, 0)
	for i := 0; i < len(listIndexes); i++ {
		indexed = append(indexed, listIndexes[i].Name)
	}

	for i := 0; i < len(keyIndexes); i++ {
		if !utilities.StringIntArray(keyIndexes[i], indexed) {
			return true
		}
	}

	return false
}

func (repo *ProductMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
