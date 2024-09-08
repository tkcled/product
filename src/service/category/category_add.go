package service_category

import (
	"context"
	"fmt"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	}

	_, err = collection.Category().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}
	return
}
