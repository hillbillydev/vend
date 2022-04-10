package vend

import (
	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Price    int       `json:"price"`
	Discount int       `json:"discount,omitempty"`

	Quantity int
}

type LineItem struct {
	ProductID uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     int
}

//ProductsTable SalesTable LineItemsTable UserTable SalesTable 1-*
//LineItemsTable SalesTable *-1 UserTable LineItemsTable 1-* ProductsTable
type Sale struct {
	Items      []LineItem
	TotalPrice int
	Quantity   int `json:"quantity"`
	UserID     uuid.UUID
}

type Storage interface {
	GetProducts() ([]Product, error)
	CreateProduct(name string, price int) (Product, error)
	Sale(sales []LineItem) ([]Product, error)
}

func TotalPrice(products []Product) int {
	var result int
	for _, p := range products {
		result += p.Price
	}
	return result
}

func ApplyDiscountOnProducts(products []Product, dollarDiscount int) []Product {
	if dollarDiscount == 0 {
		return products
	}
	var result []Product

	discount := dollarDiscount * 100 / len(products)

	for _, p := range products {
		p.Price -= discount
		p.Discount = discount
		result = append(result, p)
	}

	return result
}
