package model

import (
	"time"

	graph_model "hshelby-tkcled-product/src/graph/generated/model"
)

type Product struct {
	ID string `json:"id" bson:"_id"`

	Name        string  `json:"name" bson:"name"`
	Image       string  `json:"image" bson:"image"`
	Description string  `json:"description" bson:"description"`
	Code        string  `json:"code" bson:"code"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	CatalogLink string  `json:"catalog_link" bson:"catalog_link"`
	CategoryID  string  `json:"category_id" bson:"category_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (p *Product) ConvertToModelGraph() *graph_model.Product {
	data := graph_model.Product{
		ID: p.ID,

		Name:        p.Name,
		Image:       p.Image,
		Description: p.Description,
		Code:        p.Code,
		UnitPrice:   p.UnitPrice,
		CatalogLink: p.CatalogLink,
		Category: graph_model.Category{
			ID: p.CategoryID,
		},
	}

	return &data
}
