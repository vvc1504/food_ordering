package spec

type ProductID string

// ProductInit defines the data required to initialize a product.
type ProductInit struct {
	ID       ProductID
	Name     string
	Price    float64
	Category string
}

type Product interface {
	ID() ProductID
	Name() string
	Price() float64
	Category() string
}

type ProductManagerStatusCode StatusCode

const (
	ProductManagerNone            ProductManagerStatusCode = 0
	ProductManagerProductNotFound ProductManagerStatusCode = 1
	ProductManagerDuplicateID     ProductManagerStatusCode = 2
)

type ErrorProductManagerStatus Status

type ProductManager interface {
	NewProducts(Inits []ProductInit) (CreatedRefs <-chan []Product, St <-chan []ErrorProductManagerStatus)
	GetProductList() (Products <-chan []Product, St <-chan []ErrorProductManagerStatus)
	GetProduct(ID string) (Product <-chan Product, St <-chan ErrorProductManagerStatus)
}

// ProductObj represents a food item available for order implementation.
type ProductObj struct {
	ID_       ProductID `json:"id"`
	Name_     string    `json:"name"`
	Price_    float64   `json:"price"`
	Category_ string    `json:"category"`
}

var _ Product = (*ProductObj)(nil)

func (p *ProductObj) ID() ProductID {
	return p.ID_
}

func (p *ProductObj) Name() string {
	return p.Name_
}

func (p *ProductObj) Price() float64 {
	return p.Price_
}

func (p *ProductObj) Category() string {
	return p.Category_
}

func NewProduct(Id ProductID, Name string, Price float64, Category string) Product {
	if Id == "" || Name == "" || Price == 0 || Category == "" {
		return nil
	}
	return &ProductObj{
		ID_:       Id,
		Name_:     Name,
		Price_:    Price,
		Category_: Category,
	}
}
