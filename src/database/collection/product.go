package collection

import (
	"context"
	"fmt"
	"sync"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"

	"hshelby-tkcled-product/config"
	"hshelby-tkcled-product/src/utilities"

	"hshelby-tkcled-product/mongo"
)

const ProductCollection = "product"

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

	indexes := []mongoDriver.IndexModel{}

	if !repo.needIndex(col) {
		return
	}

	col.Indexes().CreateMany(context.Background(), indexes)
}

func (repo *ProductMongoCollection) needIndex(col *mongoDriver.Collection) bool {
	keyIndexes := []string{}

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
