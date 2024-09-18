package service_product

import (
	"context"
	"fmt"
	"log"
	"time"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAddCommand struct {
	Name        string
	Image       string
	Description string
	Code        string
	UnitPrice   float64
	CatalogLink string
	CategoryID  string
}

func (c *ProductAddCommand) Valid() error {
	if c.Code == "" {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, "miss code")
		return fmt.Errorf(codeErr)
	}

	return nil
}

func ProductAdd(ctx context.Context, c *ProductAddCommand) (result *model.Product, err error) {
	if err = c.Valid(); err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.Invalid)
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["code"] = c.Code

	cnt, err := collection.Product().Collection().CountDocuments(ctx, condition)
	if err == nil && cnt > 0 {
		log.Println("ProductAdd", map[string]interface{}{"command: ": c}, err)
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}

	if cnt > 0 {
		codeErr := fmt.Sprintf("%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.ProductCodeExist)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model.Product{
		ID: primitive.NewObjectID().Hex(),

		Name:        c.Name,
		Image:       c.Image,
		Description: c.Description,
		Code:        c.Code,
		UnitPrice:   c.UnitPrice,
		CatalogLink: c.CatalogLink,
		CategoryID:  c.CategoryID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.Product().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := fmt.Sprintf("%s-%s-%s-%s", src_const.ServiceErr_Product, src_const.ElementErr_Product, src_const.InternalError, err)
		return nil, fmt.Errorf(codeErr)
	}
	return
}
