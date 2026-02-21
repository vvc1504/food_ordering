package repository

import (
	"fmt"

	"github.com/vvc1504/food_ordering/cmd/utils"
	"github.com/vvc1504/food_ordering/pkg/spec"
)

// Repository defines the interface for data persistence.
type Repository interface {
	ProductManager() (Ref <-chan spec.ProductManager)
	OrderManager() (Ref <-chan spec.OrderManager)
}

type inMemoryRepository struct {
	ProductManagerRef_ spec.ProductManager
	OrderManagerRef_   spec.OrderManager
}

// NewInMemoryRepository creates a new in-memory implementation of the Repository.
func NewInMemoryRepository(PM spec.ProductManager, OM spec.OrderManager) (Repo Repository) {
	repo := &inMemoryRepository{
		ProductManagerRef_: PM,
		OrderManagerRef_:   OM,
	}
	return repo
}

func (r *inMemoryRepository) ProductManager() (Ref <-chan spec.ProductManager) {
	return utils.Call(func() spec.ProductManager {
		return r.ProductManagerRef_
	})
}

func (r *inMemoryRepository) OrderManager() (Ref <-chan spec.OrderManager) {
	return utils.Call(func() spec.OrderManager {
		return r.OrderManagerRef_
	})
}

func (r *inMemoryRepository) SeedData() {
	// Seeding some initial products as per the demo server style
	pm := <-r.ProductManager()
	productInits := make([]spec.ProductInit, 0)
	productInits = append(productInits, spec.ProductInit{
		ID:       spec.ProductID("1"),
		Name:     "Chicken Waffle",
		Price:    12.99,
		Category: "Waffle",
	})
	productInits = append(productInits, spec.ProductInit{
		ID:       spec.ProductID("2"),
		Name:     "Beef Burger",
		Price:    15.50,
		Category: "Burger",
	})
	productInits = append(productInits, spec.ProductInit{
		ID:       spec.ProductID("3"),
		Name:     "Veggie Pizza",
		Price:    18.00,
		Category: "Pizza",
	})
	productInits = append(productInits, spec.ProductInit{
		ID:       spec.ProductID("4"),
		Name:     "Classic Fries",
		Price:    4.50,
		Category: "Sides",
	})
	productInits = append(productInits, spec.ProductInit{
		ID:       spec.ProductID("5"),
		Name:     "Iced Tea",
		Price:    3.00,
		Category: "Beverage",
	})
	productsListChan, errListChan := pm.NewProducts(productInits)
	_ = <-productsListChan
	errList := <-errListChan
	if errList != nil {
		fmt.Println("Error seeding products:", errList)
	}
}
