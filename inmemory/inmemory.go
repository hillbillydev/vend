package inmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/vend/vend"
)

type InMemory struct {
	mtx      *sync.Mutex
	products map[string]vend.Product
}

// NewInMemoryInMemory test
func NewInMemory() *InMemory {
	return &InMemory{
		mtx:      &sync.Mutex{},
		products: make(map[string]vend.Product),
	}
}

func (m *InMemory) GetProducts() ([]vend.Product, error) {
	var result []vend.Product
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, p := range m.products {
		result = append(result, p)
	}

	return result, nil
}

func (m *InMemory) CreateProduct(name string, price int) (vend.Product, error) {
	id := uuid.New()
	p := vend.Product{
		ID:    id,
		Name:  name,
		Price: price,
	}

	m.mtx.Lock()
	m.products[id.String()] = p
	m.mtx.Unlock()

	return p, nil
}

func (m *InMemory) Sale(sales []vend.LineItem) ([]vend.Product, error) {
	var result []vend.Product

	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, sale := range sales {
		p, ok := m.products[sale.ProductID.String()]
		if !ok {
			return nil, fmt.Errorf("product with id of %s does not exists in the database", sale.ProductID.String())
		}

		for i := sale.Quantity; i > 0; i-- {
			result = append(result, p)
		}
	}

	return result, nil
}
