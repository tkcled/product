// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph_model

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryAdd struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

type CategoryDelete struct {
	CategoryID string `json:"category_id"`
}

type CategoryPagination struct {
	Rows   []Product  `json:"rows"`
	Paging Pagination `json:"paging"`
}

type CategoryUpdate struct {
	CategoryID  string  `json:"category_id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type DefaultResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Mutation struct {
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalPage   int `json:"total_page"`
	Total       int `json:"total"`
}

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
	UnitPrice   float64  `json:"unit_price"`
	CatalogLink string   `json:"catalog_link"`
	Category    Category `json:"category"`
}

type ProductAdd struct {
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	UnitPrice   float64 `json:"unit_price"`
	CatalogLink string  `json:"catalog_link"`
	CategoryID  string  `json:"category_id"`
}

type ProductDelete struct {
	ID string `json:"id"`
}

type ProductPagination struct {
	Rows   []Product  `json:"rows"`
	Paging Pagination `json:"paging"`
}

type ProductUpdate struct {
	ID          string   `json:"id"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Code        *string  `json:"code,omitempty"`
	UnitPrice   *float64 `json:"unit_price,omitempty"`
	CatalogLink *string  `json:"catalog_link,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

type Query struct {
}
