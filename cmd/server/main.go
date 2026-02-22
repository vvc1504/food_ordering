package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vvc1504/food_ordering/pkg/handlers"
	"github.com/vvc1504/food_ordering/pkg/middleware"
	"github.com/vvc1504/food_ordering/pkg/spec"
	"github.com/vvc1504/food_ordering/pkg/spec/impl"
	"github.com/vvc1504/food_ordering/pkg/spec/impl/coupon"
)

func main() {
	// Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Flags
	port := flag.String("port", "8080", "Server port")
	couponFilesStr := flag.String("coupon-files", "requirements/couponbase1.gz,requirements/couponbase2.gz,requirements/couponbase3.gz", "Comma-separated list of coupon gzip files")
	flag.Parse()

	log.Info().Msg("Starting Food Ordering API server...")

	// We will use an atomic flag to mark when the API is ready
	var isReady atomic.Bool
	isReady.Store(false)

	// API Wrapper to return 503 when not ready
	apiWrapper := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !isReady.Load() {
				http.Error(w, "Server is initializing coupons, please wait...", http.StatusServiceUnavailable)
				return
			}
			next(w, r)
		}
	}

	// Setup Router and Middleware early so frontend can be served
	mux := http.NewServeMux()

	// 1. Initialize Handlers early but with nil dependencies (they won't be called until isReady=true)
	// Actually, we can wrap the mux definition to use a dynamic handler or just defer handler execution
	var h *handlers.Handler

	// Public routes
	mux.HandleFunc("/product", middleware.CORS(middleware.Logger(apiWrapper(func(w http.ResponseWriter, r *http.Request) { h.ListProducts(w, r) }))))
	mux.HandleFunc("/product/", middleware.CORS(middleware.Logger(apiWrapper(func(w http.ResponseWriter, r *http.Request) { h.GetProduct(w, r) }))))
	mux.HandleFunc("/orders", middleware.CORS(middleware.Logger(apiWrapper(func(w http.ResponseWriter, r *http.Request) { h.ListOrders(w, r) }))))

	// Protected routes
	mux.HandleFunc("/order", middleware.CORS(middleware.Logger(middleware.APIKeyAuth("apitest", apiWrapper(func(w http.ResponseWriter, r *http.Request) { h.PlaceOrder(w, r) })))))

	// Serve Frontend Static Files
	mux.Handle("/", http.FileServer(http.Dir("web/dist")))

	// Start server in the background
	addr := fmt.Sprintf(":%s", *port)
	log.Info().Msgf("Server listening on %s", addr)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// --- Heavy Initialization Phase ---

	// 2. Initialize Coupon Validator (Sequential Streaming Indexing)
	couponFiles := strings.Split(*couponFilesStr, ",")

	log.Info().Msg("Indexing coupons (this may take a moment)...")
	couponValidator, err := coupon.NewValidator(couponFiles)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize coupon validator")
	}
	log.Info().Msg("Coupon indexing complete.")

	// 3. Initialize Managers
	productMgr := impl.NewProductManagerRef()
	prdInitStChan := <-productMgr.(impl.ProductManagerInternal).ProductManagerInit()
	_ = prdInitStChan

	orderMgr := impl.NewOrderManagerRef()
	initStChan := <-orderMgr.(impl.OrderManagerInternal).OrderManagerInit(productMgr, couponValidator)
	_ = initStChan

	// 4. Initialize Products from "Config" (Seed data)
	productInits := []spec.ProductInit{
		{ID: spec.ProductID("1"), Name: "Chicken Waffle", Price: 12.99, Category: "Waffle"},
		{ID: spec.ProductID("2"), Name: "Beef Burger", Price: 15.50, Category: "Burger"},
		{ID: spec.ProductID("3"), Name: "Veggie Pizza", Price: 18.00, Category: "Pizza"},
		{ID: spec.ProductID("4"), Name: "Classic Fries", Price: 4.50, Category: "Sides"},
		{ID: spec.ProductID("5"), Name: "Iced Tea", Price: 3.00, Category: "Beverage"},
	}

	log.Info().Msg("Initializing products aggregate...")
	initRefsChan, initSeedingStChan := productMgr.NewProducts(productInits)
	<-initRefsChan // Wait for refs
	initSeedingSts := <-initSeedingStChan

	for _, st := range initSeedingSts {
		if st != nil && st.Code() != 0 {
			log.Fatal().Interface("status", st).Msg("Failed to initialize products")
		}
	}
	log.Info().Msg("Products aggregate initialized.")

	// 5. Assign dependencies and mark AS READY
	h = handlers.NewHandler(productMgr, orderMgr)
	isReady.Store(true)

	log.Info().Msg("Server is fully ready. API endpoints are now accepting traffic.")

	// Block main thread to keep server running
	select {}
}
