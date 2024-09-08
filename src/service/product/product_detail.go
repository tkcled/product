package service_product

import (
	"context"
	"fmt"
	"log"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductDetailCommand struct {
	ProductID string
}

func (c *ProductDetailCommand) Valid() error {
	if c.ProductID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, "miss product id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ProductDetail(ctx context.Context, c *ProductDetailCommand) (result *model.Product, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Product{}
	err = collection.Product().Collection().FindOne(ctx, bson.M{"id": c.ProductID}).Decode(result)
	if err != nil {
		log.Println("ProductDetail", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
