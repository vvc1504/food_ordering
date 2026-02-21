package impl

import (
	"fmt"

	"github.com/vvc1504/food_ordering/cmd/utils"
	"github.com/vvc1504/food_ordering/pkg/spec"
)

// orderManager implements the spec.OrderManager interface.
type OrderManager struct {
	ProductManagerRef_ spec.ProductManager
	Orders             map[string]spec.Order
	OrderKeys          []string
	CouponValidator    spec.CouponValidator
}

type OrderManagerInternal interface {
	OrderManagerInit(ProductManagerRef spec.ProductManager, CV spec.CouponValidator) (St <-chan spec.ErrorOrderManagerStatus)
}

func NewOrderManagerRef() (Ref spec.OrderManager) {
	return &OrderManager{}
}

func (p *OrderManager) ProductManager() (Ref <-chan spec.ProductManager) {
	return utils.Call(func() spec.ProductManager {
		return p.ProductManagerRef_
	})
}

func (p *OrderManager) HasOrderKey(ID string) (Has <-chan bool) {
	_, ok := p.Orders[ID]
	return utils.Call(func() bool {
		return ok
	})
}

func (m *OrderManager) OrderManagerInit(ProductManagerRef spec.ProductManager, CV spec.CouponValidator) (St <-chan spec.ErrorOrderManagerStatus) {
	return utils.Call(func() spec.ErrorOrderManagerStatus {
		m.CouponValidator = CV
		m.ProductManagerRef_ = ProductManagerRef
		m.Orders = make(map[string]spec.Order, 0)
		m.OrderKeys = make([]string, 0)
		return nil
	})
}

// PlaceOrder implements [spec.OrderManager].
func (m *OrderManager) PlaceOrder(Reqs []spec.OrderReq) (Order <-chan []spec.Order, St <-chan []spec.ErrorOrderManagerStatus) {
	created := make([]spec.Order, 0, len(Reqs))
	statuses := make([]spec.ErrorOrderManagerStatus, 0, len(Reqs))
	pm := <-m.ProductManager()
	for _, req := range Reqs {
		if len(req.Items) == 0 {
			statuses = append(statuses, (spec.ErrorOrderManagerStatus)(spec.NewStatus(spec.StatusCode(spec.OrderManagerInvalidInput), "order must contain at least one item")))
			continue
		}
		if req.CouponCode != "" {
			valid, err := m.CouponValidator.Validate(req.CouponCode)
			if err != nil || !valid {
				statuses = append(statuses, (spec.ErrorOrderManagerStatus)(spec.NewStatus(spec.StatusCode(spec.OrderManagerInvalidCoupon), "invalid coupon code")))
				continue
			}
		}
		orderID := spec.OrderID(utils.GenerateID())
		order := &spec.OrderObj{
			ID_:    orderID,
			Items_: req.Items,
		}
		m.Orders[string(orderID)] = order
		m.OrderKeys = append(m.OrderKeys, string(orderID))
		created = append(created, order)
		statuses = append(statuses, nil)

		for _, item := range req.Items {
			prdChan, prdStChan := pm.GetProduct(string(item.ProductID))
			prd := <-prdChan
			prdSt := <-prdStChan
			if prdSt != nil && prdSt.Code() != 0 {
				statuses = append(statuses, (spec.ErrorOrderManagerStatus)(spec.NewStatus(spec.StatusCode(spec.OrderManagerProductMissing), fmt.Sprintf("product %s not found", item.ProductID))))
				continue
			}
			order.Products_ = append(order.Products_, prd)
		}
	}
	return utils.Call2(func() ([]spec.Order, []spec.ErrorOrderManagerStatus) {
		return created, statuses
	})
}

var _ OrderManagerInternal = (*OrderManager)(nil)
var _ spec.OrderManager = (*OrderManager)(nil)
