package model

import (
	"time"

	graph_model "hshelby-tkcled-product/src/graph/generated/model"
)

type Category struct {
	ID string `json:"id" bson:"_id"`

	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	ParentID    string      `json:"parent_id" bson:"parent_id"`
	Children    *[]Category `json:"children,omitempty"`
	Seq         int         `json:"seq" bson:"seq"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (c *Category) ConvertToModelGraph() *graph_model.Category {
	data := graph_model.Category{
		ID: c.ID,

		Name:        c.Name,
		Description: c.Description,
		Parent: &graph_model.Category{
			ID: c.ParentID,
		},
		Seq: c.Seq,
	}

	if c.Children != nil {
		data.Children = func() []graph_model.Category {
			items := make([]graph_model.Category, 0)
			for _, ele := range *c.Children {
				items = append(items, *ele.ConvertToModelGraph())
			}
			return items
		}()
	}

	return &data
}
