package service_product

import (
	"context"
	"fmt"
	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductDeleteCommand struct {
	ProductID string
}

func (c *ProductDeleteCommand) Valid() error {
	if c.ProductID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, "miss product id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ProductDelete(ctx context.Context, c *ProductDeleteCommand) (result *model.Product, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Product{}
	err = collection.Product().Collection().FindOne(ctx, bson.M{"id": c.ProductID}).Decode(result)
	if err != nil {
		log.Println("ProductDelete", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["_id"] = c.ProductID

	_, err = collection.Product().Collection().DeleteOne(ctx, condition)
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
