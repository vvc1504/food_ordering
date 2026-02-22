package impl

import (
	"sync"

	"github.com/vvc1504/food_ordering/cmd/utils"
	"github.com/vvc1504/food_ordering/pkg/spec"
)

type ProductManagerInternal interface {
	ProductManagerInit() (St <-chan spec.ErrorProductManagerStatus)
}

func NewProductManagerRef() (Ref spec.ProductManager) {
	return &ProductManager{
		Products:    make(map[string]spec.Product),
		ProductKeys: []string{},
	}
}

// productManager implements the spec.ProductManager interface.
type ProductManager struct {
	mu          sync.RWMutex
	Products    map[string]spec.Product
	ProductKeys []string
}

func (p *ProductManager) HasProductKey(ID string) (Has <-chan bool) {
	return utils.Call(func() bool {
		p.mu.RLock()
		defer p.mu.RUnlock()
		_, ok := p.Products[ID]
		return ok
	})
}

// GetProduct implements [spec.ProductManager].
func (p *ProductManager) GetProduct(ID string) (Product <-chan spec.Product, St <-chan spec.ErrorProductManagerStatus) {
	return utils.Call2(func() (spec.Product, spec.ErrorProductManagerStatus) {
		p.mu.RLock()
		defer p.mu.RUnlock()
		prd, ok := p.Products[ID]
		if !ok {
			return nil, (spec.ErrorProductManagerStatus)(spec.NewStatus(spec.StatusCode(spec.ProductManagerProductNotFound), "product not found"))
		}
		return prd, nil
	})
}

// GetProductList implements [spec.ProductManager].
func (p *ProductManager) GetProductList() (Products <-chan []spec.Product, St <-chan []spec.ErrorProductManagerStatus) {
	return utils.Call2(func() ([]spec.Product, []spec.ErrorProductManagerStatus) {
		p.mu.RLock()
		defer p.mu.RUnlock()

		products := make([]spec.Product, 0, len(p.ProductKeys))
		statuses := make([]spec.ErrorProductManagerStatus, 0, len(p.ProductKeys))
		for _, pkey := range p.ProductKeys {
			prd, ok := p.Products[pkey]
			if !ok {
				statuses = append(statuses, (spec.ErrorProductManagerStatus)(spec.NewStatus(spec.StatusCode(spec.ProductManagerProductNotFound), "product not found")))
				continue
			} else {
				products = append(products, prd)
				statuses = append(statuses, nil)
			}
		}
		return products, statuses
	})
}

// NewProducts implements [spec.ProductManager].
func (p *ProductManager) NewProducts(Inits []spec.ProductInit) (CreatedRefs <-chan []spec.Product, St <-chan []spec.ErrorProductManagerStatus) {
	return utils.Call2(func() ([]spec.Product, []spec.ErrorProductManagerStatus) {
		p.mu.Lock()
		defer p.mu.Unlock()

		created := make([]spec.Product, 0, len(Inits))
		statuses := make([]spec.ErrorProductManagerStatus, 0, len(Inits))
		for _, init := range Inits {
			_, ok := p.Products[string(init.ID)]
			if ok {
				statuses = append(statuses, (spec.ErrorProductManagerStatus)(spec.NewStatus(spec.StatusCode(spec.ProductManagerDuplicateID), "product already exists")))
				continue
			}
			prd := spec.NewProduct(init.ID, init.Name, init.Price, init.Category)
			p.Products[string(init.ID)] = prd
			p.ProductKeys = append(p.ProductKeys, string(init.ID))
			created = append(created, prd)
			statuses = append(statuses, nil)
		}
		return created, statuses
	})
}

// ProductManagerInit implements [ProductManagerInternal].
func (p *ProductManager) ProductManagerInit() (St <-chan spec.ErrorProductManagerStatus) {
	return utils.Call(func() spec.ErrorProductManagerStatus {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.Products = make(map[string]spec.Product, 0)
		p.ProductKeys = make([]string, 0)
		return nil
	})
}

var _ spec.ProductManager = (*ProductManager)(nil)
var _ ProductManagerInternal = (*ProductManager)(nil)
