package service_product

import (
	"context"
	"fmt"
	src_const "hshelby-tkcled-product/src/const"
	"hshelby-tkcled-product/src/database/collection"
	"hshelby-tkcled-product/src/database/model"
	"log"
	"strings"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ProductPaginationCommand struct {
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	OrderBy string                 `json:"order_by"`
	Search  map[string]interface{} `json:"search"`
}

func (c *ProductPaginationCommand) Valid() error {
	if c.Page < 1 {
		c.Page = 1
	}

	if c.Limit < 1 {
		c.Limit = 10
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.Invalid
		return fmt.Errorf(codeErr)
	}
	return nil
}

func ProductPagination(ctx context.Context, c *ProductPaginationCommand) (total int, results []model.Product, err error) {
	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.Invalid
		return 0, nil, fmt.Errorf(codeErr)
	}
	condition := make(map[string]interface{})

	if name, ok := c.Search["name"]; ok {
		condition["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	if categoryID, ok := c.Search["category_id"]; ok {
		condition["$or"] = []bson.M{
			{"category_id": bson.M{"$regex": categoryID, "$options": "i"}}, // Condition 1
			{"parent_category_id": bson.M{"$regex": categoryID, "$options": "i"}},
		}
	}

	if code, ok := c.Search["code"]; ok {
		condition["code"] = code
	}

	objOrderBy := bson.M{}
	if c.OrderBy != "" {
		value := src_const.ASC
		if strings.HasPrefix(c.OrderBy, "-") {
			value = src_const.DESC
			c.OrderBy = strings.TrimPrefix(c.OrderBy, "-")
		}

		objOrderBy = bson.M{c.OrderBy: value}
	}

	//Default order by updated_at | new -> old
	if c.OrderBy == "" {
		objOrderBy = bson.M{"name": src_const.DESC}
	}

	matchStage := bson.D{{Key: "$match", Value: condition}}

	facectStage := bson.D{{
		Key: "$facet",

		Value: bson.M{
			"rows": bson.A{
				bson.M{"$skip": (c.Page - 1) * c.Limit},
				bson.M{"$limit": c.Limit},
			},
			"total": bson.A{
				bson.M{"$count": "count"},
			},
		},
	}}

	fmt.Println(condition)

	sortStage := bson.D{{Key: "$sort", Value: objOrderBy}}

	pipeline := mongoDriver.Pipeline{
		matchStage,
		sortStage,
		facectStage,
	}

	cur, err := collection.Product().Collection().Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("ProductPagination", err)
		codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.InternalError
		return 0, nil, fmt.Errorf(codeErr)
	}

	var listOrder bson.M
	for cur.Next(ctx) {
		err := cur.Decode(&listOrder)
		if err != nil {
			log.Println("ProductPagination", err)
			codeErr := src_const.ServiceErr_Product + src_const.ServiceErr_Product + src_const.InternalError
			return 0, nil, fmt.Errorf(codeErr)
		}
	}

	// Extract the total count and rows from the result
	products := make([]model.Product, 0)

	if len(listOrder["total"].(bson.A)) > 0 {
		total = int(listOrder["total"].(bson.A)[0].(bson.M)["count"].(int32))
		rows := listOrder["rows"].(bson.A)

		var product model.Product

		for _, rawProduct := range rows {
			productBSON, err := bson.Marshal(rawProduct)
			if err != nil {
				log.Println("ProductPagination", err)
				codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.InternalError
				return 0, nil, fmt.Errorf(codeErr)
			}

			err = bson.Unmarshal(productBSON, &product)
			if err != nil {
				log.Println("ProductPagination", err)
				codeErr := src_const.ServiceErr_Product + src_const.ElementErr_Product + src_const.InternalError
				return 0, nil, fmt.Errorf(codeErr)
			}

			products = append(products, product)
		}
	}

	return int(total), products, nil
}
