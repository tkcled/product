package service_category

import (
	"context"
	"fmt"
	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type CategoryDeleteCommand struct {
	CategoryID string
}

func (c *CategoryDeleteCommand) Valid() error {
	if c.CategoryID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, "miss Category id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func CategoryDelete(ctx context.Context, c *CategoryDeleteCommand) (result *model.Category, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Category{}
	err = collection.Category().Collection().FindOne(ctx, bson.M{"id": c.CategoryID}).Decode(result)
	if err != nil {
		log.Println("CategoryDelete", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["_id"] = c.CategoryID

	_, err = collection.Category().Collection().DeleteOne(ctx, condition)
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
