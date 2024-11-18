package service_product

import (
	"context"
	"fmt"
	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	"hshelby-tkcled-product/src/database/model"
	"log"
)

func SpecialProduct(ctx context.Context) (results []model.Product, err error) {
	condition := make(map[string]interface{})
	condition["is_special"] = true
	cur, err := collection.Product().Collection().Find(ctx, condition)
	if err != nil {
		log.Println("err", err)
		codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}

	results = make([]model.Product, 0)
	err = cur.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
