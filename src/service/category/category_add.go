package service_category

import (
	"context"
	"fmt"
	"log"
	"time"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryAddCommand struct {
	Name        string
	Description string
	ParentID    string
}

func CategoryAdd(ctx context.Context, c *CategoryAddCommand) (result *model.Category, err error) {
	result = &model.Category{
		ID: primitive.NewObjectID().Hex(),

		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if c.ParentID != "" {
		parent := &model.Category{}
		err = collection.Category().Collection().FindOne(ctx, bson.M{"_id": c.ParentID}).Decode(parent)
		if err != nil {
			log.Println("CategoryAdd", map[string]interface{}{"command: ": c}, err)
			return &model.Category{}, nil
		}
		result.Seq = parent.Seq
	}

	if c.ParentID == "" {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"seq", -1}})
		findOptions.SetLimit(1)

		cursor, err := collection.Category().Collection().Find(ctx, bson.M{}, findOptions)
		if err != nil {
			log.Println("CategoryAdd", map[string]interface{}{"command: ": c}, err)
		}

		if err == nil {
			categories := make([]model.Category, 0)
			err = cursor.All(ctx, &categories)
			if err != nil {
				log.Println("CategoryAdd", map[string]interface{}{"command: ": c}, err)
			}

			result.Seq = categories[0].Seq + 1

		}
	}

	_, err = collection.Category().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}
	return
}
