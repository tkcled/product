package model

import (
	"time"

	graph_model "hshelby-tkcled-product/src/graph/generated/model"
)

type Category struct {
	ID string `json:"id" bson:"_id"`

	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (c *Category) ConvertToModelGraph() *graph_model.Category {
	data := graph_model.Category{
		ID: c.ID,

		Name:        c.Name,
		Description: c.Description,
	}

	return &data
}
