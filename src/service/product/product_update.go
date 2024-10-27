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

type ProductUpdateCommand struct {
	ProductID string

	Name        *string
	Description *string
	Code        *string
	UnitPrice   *float64
	CatalogLink *string
	CategoryID  *string
}

func (c *ProductUpdateCommand) Valid() error {
	if c.ProductID == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, "miss product id")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ProductUpdate(ctx context.Context, c *ProductUpdateCommand) (result *model.Product, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Product{}
	err = collection.Product().Collection().FindOne(ctx, bson.M{"_id": c.ProductID}).Decode(result)
	if err != nil {
		log.Println("ProductUpdate", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	updated := make(map[string]interface{})

	if c.Name != nil && c.Name != &result.Name {
		updated["name"] = *c.Name
	}

	if c.Description != nil && c.Description != &result.Description {
		updated["description"] = *c.Description
	}

	if c.Code != nil && c.Code != &result.Code {
		// condition := make(map[string]interface{})
		// condition["code"] = *c.Code
		// cnt, err := collection.Product().Collection().CountDocuments(ctx, condition)
		// if err != nil {
		// 	log.Println("ProductUpdate", map[string]interface{}{"command: ": c}, err)
		// 	codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		// 	return nil, fmt.Errorf(codeErr)
		// }

		// if cnt > 0 {
		// 	codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.ProductCodeExist)
		// 	return nil, fmt.Errorf(codeErr)
		// }

		updated["code"] = *c.Code
	}

	if c.UnitPrice != nil && c.UnitPrice != &result.UnitPrice {
		updated["unit_price"] = *c.UnitPrice
	}

	if c.CatalogLink != nil && c.CatalogLink != &result.CatalogLink {
		updated["catalog_link"] = *c.CatalogLink
	}

	if c.CategoryID != nil && c.CategoryID != &result.CategoryID {
		updated["category_id"] = *c.CategoryID
	}

	_, err = collection.Product().Collection().UpdateByID(ctx, c.ProductID, bson.M{"$set": updated})
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	return result, nil
}
