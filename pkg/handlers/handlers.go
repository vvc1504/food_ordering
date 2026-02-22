package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vvc1504/food_ordering/pkg/spec"
)

type Handler struct {
	productMgr spec.ProductManager
	orderMgr   spec.OrderManager
}

func NewHandler(pm spec.ProductManager, om spec.OrderManager) *Handler {
	return &Handler{
		productMgr: pm,
		orderMgr:   om,
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productsChan, stsChan := h.productMgr.GetProductList()
	products := <-productsChan
	sts := <-stsChan

	for _, st := range sts {
		if st != nil && st.Code() != 0 {
			http.Error(w, st.Message(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	id := parts[2]

	productChan, stChan := h.productMgr.GetProduct(id)
	product := <-productChan
	st := <-stChan

	if st != nil && st.Code() != 0 {
		if st.Code() == 1 { // ProductNotFound
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, st.Message(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req spec.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	orderChan, stsChan := h.orderMgr.PlaceOrder([]spec.OrderReq{req})
	orders := <-orderChan
	sts := <-stsChan

	if len(sts) > 0 && sts[0] != nil && sts[0].Code() != 0 {
		st := sts[0]
		status := http.StatusInternalServerError
		if st.Code() == 1 || st.Code() == 2 || st.Code() == 3 {
			status = http.StatusUnprocessableEntity
		}
		http.Error(w, st.Message(), status)
		return
	}

	if len(orders) == 0 {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"orderId": string(orders[0].ID())})
}
