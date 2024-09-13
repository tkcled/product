package service_category

import (
	"context"
	"fmt"
	"log"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson"
)

type CategoryDetailCommand struct {
	CategoryID string
}

func (c *CategoryDetailCommand) Valid() error {
	if c.CategoryID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, "miss category id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func CategoryDetail(ctx context.Context, c *CategoryDetailCommand) (result *model.Category, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Category{}
	err = collection.Category().Collection().FindOne(ctx, bson.M{"_id": c.CategoryID}).Decode(result)
	if err != nil {
		log.Println("CategoryDetail", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Category, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
