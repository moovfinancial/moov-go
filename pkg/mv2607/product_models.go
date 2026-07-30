package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

// Product is a good or service offered by a merchant, including the product taxonomy
// category associated with it.
type Product struct {
	moov.Product

	// Category is the product taxonomy category associated with the product, if any.
	Category *ProductCategory `json:"category,omitempty"`
}

// ProductRequest is a request to create or update a product.
type ProductRequest struct {
	Title        string                          `json:"title"`
	Description  *string                         `json:"description,omitempty"`
	BasePrice    moov.AmountDecimal              `json:"basePrice"`
	Images       []moov.AssignProductImage       `json:"images,omitempty"`
	OptionGroups []moov.CreateProductOptionGroup `json:"optionGroups,omitempty"`

	// CategoryID is the ID of a product taxonomy category to associate with the product.
	CategoryID *string `json:"categoryID,omitempty"`
}

// ProductCategory is a product category from the product taxonomy.
type ProductCategory struct {
	// CategoryID is the unique identifier for the category.
	CategoryID string `json:"categoryID"`
	// Name is the short display name of the category, e.g. "Beverages".
	Name string `json:"name"`
	// FullName is the full taxonomy path name of the category,
	// e.g. "Food, Beverages & Tobacco > Beverages".
	FullName string `json:"fullName"`
	// Level is the depth of the category in the taxonomy tree (1 = top-level).
	Level int32 `json:"level"`
	// ParentID is the identifier of the parent category. Absent for top-level categories.
	ParentID *string `json:"parentID,omitempty"`
}

// ProductCategories is a list of product categories from the product taxonomy.
type ProductCategories struct {
	Categories []ProductCategory `json:"categories"`
}
