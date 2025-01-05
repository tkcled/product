package service_product

import (
	"context"
	"fmt"
	"time"

	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	model "hshelby-tkcled-product/src/database/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAddCommand struct {
	Name             string
	Description      string
	Code             string
	UnitPrice        float64
	CatalogLink      string
	CategoryID       string
	ParentCategoryID string
	Metadata         string
	IsSpecial        bool
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

	result = &model.Product{
		ID: primitive.NewObjectID().Hex(),

		Name:             c.Name,
		Description:      c.Description,
		Code:             c.Code,
		UnitPrice:        c.UnitPrice,
		CatalogLink:      c.CatalogLink,
		CategoryID:       c.CategoryID,
		ParentCategoryID: c.ParentCategoryID,
		Metadata:         c.Metadata,
		IsSpecial:        c.IsSpecial,

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
