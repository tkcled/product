package model

import (
	"strings"
	"time"

	graph_model "hshelby-tkcled-product/src/graph/generated/model"
)

type Product struct {
	ID string `json:"id" bson:"_id"`

	Name             string  `json:"name" bson:"name"`
	Image            string  `json:"image" bson:"image"`
	Description      string  `json:"description" bson:"description"`
	Code             string  `json:"code" bson:"code"`
	UnitPrice        float64 `json:"unit_price" bson:"unit_price"`
	CatalogLink      string  `json:"catalog_link" bson:"catalog_link"`
	CategoryID       string  `json:"category_id" bson:"category_id"`
	ParentCategoryID string  `json:"parent_category_id" bson:"parent_category_id"`
	Metadata         string  `json:"metadata" bson:"metadata"`
	IsSpecial        bool    `json:"is_special" bson:"is_special"`

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
		Metadata:    p.Metadata,
		IsSpecial:   p.IsSpecial,
	}

	if p.CategoryID != "" {
		listCate := strings.Split(p.CategoryID, ",")
		for _, categoryID := range listCate {
			data.Category = append(data.Category, graph_model.Category{
				ID: categoryID,
			})
		}
	}

	if p.ParentCategoryID != "" {
		listCate := strings.Split(p.ParentCategoryID, ",")
		for _, categoryID := range listCate {
			data.ParentCategory = append(data.ParentCategory, graph_model.Category{
				ID: categoryID,
			})
		}
	}

	return &data
}
