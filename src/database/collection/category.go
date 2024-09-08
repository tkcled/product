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

const CategoryCollection = "category"

var (
	_categoryCollection        *CategoryMongoCollection
	loadCategoryRepositoryOnce sync.Once
)

type CategoryMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadCategoryCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadCategoryRepositoryOnce.Do(func() {
		_categoryCollection, err = NewCategoryMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Category() *CategoryMongoCollection {
	if _categoryCollection == nil {
		panic("database: like category collection is not initiated")
	}
	return _categoryCollection
}

func NewCategoryMongoCollection(client *mongo.MongoDB, databaseName string) (*CategoryMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewCategoryMongoCollection] client nil pointer")
	}
	repo := &CategoryMongoCollection{
		client:         client,
		collectionName: CategoryCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	repo.SetIndex()
	return repo, nil
}

func (repo *CategoryMongoCollection) SetIndex() {
	col := repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)

	indexes := []mongoDriver.IndexModel{}

	if !repo.needIndex(col) {
		return
	}

	col.Indexes().CreateMany(context.Background(), indexes)
}

func (repo *CategoryMongoCollection) needIndex(col *mongoDriver.Collection) bool {
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

func (repo *CategoryMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
