package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	couponDir := flag.String("coupon-dir", "requirements", "Directory containing couponbaseX.gz files")
	flag.Parse()

	log.Info().Msg("Starting Food Ordering API server...")

	// 1. Initialize Repository

	// 2. Initialize Coupon Validator (Sequential Streaming Indexing)
	couponFiles := []string{
		filepath.Join(*couponDir, "couponbase1.gz"),
		filepath.Join(*couponDir, "couponbase2.gz"),
		filepath.Join(*couponDir, "couponbase3.gz"),
	}

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

	// 5. Initialize Handlers
	h := handlers.NewHandler(productMgr, orderMgr)

	// 5. Setup Router and Middleware
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/product", middleware.Logger(h.ListProducts))
	mux.HandleFunc("/product/", middleware.Logger(h.GetProduct))

	// Protected routes
	mux.HandleFunc("/order", middleware.Logger(middleware.APIKeyAuth("apitest", h.PlaceOrder)))

	// Start server
	addr := fmt.Sprintf(":%s", *port)
	log.Info().Msgf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
