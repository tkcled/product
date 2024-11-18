package service_category

import (
	"context"
	"fmt"
	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	"hshelby-tkcled-product/src/database/model"
	"log"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryListCommand struct {
	ParentID string
}

func ListCategory(ctx context.Context, c CategoryListCommand) (results []model.Category, err error) {
	log.Println("[service_category.ListCategory] start")
	defer func() {
		log.Println("[service_category.ListCategory] end error", err)
	}()

	findOptions := options.Find().SetSort(bson.D{{"seq", 1}})
	cur, err := collection.Category().Collection().Find(ctx, bson.M{"parent_id": c.ParentID}, findOptions)
	if err != nil {
		codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Category + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}

	results = make([]model.Category, 0)
	err = cur.All(ctx, &results)
	if err != nil {
		return results, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}
