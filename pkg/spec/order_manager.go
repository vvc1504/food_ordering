package spec

type OrderID string

// OrderItem represents a single item in an order.
type OrderItem struct {
	ProductID ProductID `json:"productId"`
	Quantity  int       `json:"quantity"`
}

type Order interface {
	ID() (ID OrderID)
	Items() (Items []OrderItem)
	Products() (Products []Product)
}

// OrderObj represents a placed order implementation.
type OrderObj struct {
	ID_       OrderID     `json:"id"`
	Items_    []OrderItem `json:"items"`
	Products_ []Product   `json:"products"`
}

var _ Order = (*OrderObj)(nil)

func (o *OrderObj) ID() (ID OrderID)               { return o.ID_ }
func (o *OrderObj) Items() (Items []OrderItem)     { return o.Items_ }
func (o *OrderObj) Products() (Products []Product) { return o.Products_ }

// OrderReq represents the request body for placing a new order.
type OrderReq struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

type OrderManagerStatusCode StatusCode

const (
	OrderManagerNone           OrderManagerStatusCode = 0
	OrderManagerInvalidInput   OrderManagerStatusCode = 1
	OrderManagerInvalidCoupon  OrderManagerStatusCode = 2
	OrderManagerProductMissing OrderManagerStatusCode = 3
)

type ErrorOrderManagerStatus Status

type OrderManager interface {
	PlaceOrder(Reqs []OrderReq) (Order <-chan []Order, St <-chan []ErrorOrderManagerStatus)
}
