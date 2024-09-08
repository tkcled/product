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

type CategoryUpdateCommand struct {
	CategoryID string

	Name        *string
	Description *string
}

func (c *CategoryUpdateCommand) Valid() error {
	if c.CategoryID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, "miss product id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func CategoryUpdate(ctx context.Context, c *CategoryUpdateCommand) (result *model.Category, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Category{}
	err = collection.Product().Collection().FindOne(ctx, bson.M{"_id": c.CategoryID}).Decode(result)
	if err != nil {
		log.Println("CategoryUpdate", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	updated := make(map[string]interface{})

	if c.Name != nil {
		updated["name"] = *c.Name
	}

	if c.Name != nil {
		updated["name"] = *c.Name
	}

	if c.Description != nil {
		updated["description"] = *c.Description
	}

	_, err = collection.Category().Collection().UpdateByID(ctx, c.CategoryID, bson.M{"$set": updated})
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
