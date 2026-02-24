# Food Ordering Platform

This repository contains a full-stack food ordering application with a Go backend and a React frontend. The project implements a robust API server based on the OpenAPI 3.1 specification and a high-performance coupon validation engine.

## 🚀 Key Features

- **Full-Stack Implementation**: Modern React frontend and high-performance Go backend.
- **Clean Architecture**: Organized into logical layers separating API handlers, business logic, and data.
- **Fast Coupon Lookup**: Custom validation engine that streams and indexes ~2GB of compressed coupon data with low memory usage.
- **Flexible Data Model**: Implementation follows a robust manager pattern with non-blocking async logic.
- **Bonus Discounts**:
  - `HAPPYHOURS`: Applies an 18% discount to the order total.
  - `BUYGETONE`: Gives the lowest-priced item for free (for every 2 items bought).
- **Responsive Design**: Mobile-first UI optimized for both 375px and 1440px widths.
- **Dockerized**: Easy deployment using a multi-stage Docker build.

## 🏗️ Technical Overview

The backend is built with a focus on production-grade extensibility:
- **Service Managers**: `ProductManager` and `OrderManager` handle the core logic.
- **Async Pattern**: Operations return results via Go channels for non-blocking communication.
- **In-Memory Storage**: Fast, thread-safe in-memory maps for products and orders.
- **Security**: Placing orders requires an `api_key: apitest` header.

### High-Performance Coupon Validation
The challenge involved validating millions of coupons across three large `.gz` files:
1. **Streaming**: Files are read directly via `gzip` streams, never touching the disk as uncompressed data.
2. **Logic**: A coupon is valid only if it appears in at least two of the source files.
3. **Efficiency**: Temporary indexing maps are cleared immediately after the final valid set is built, keeping memory footprint low.

---

## 🛠️ Getting Started (All Platforms)

### Option 1: Docker (Recommended - Linux, macOS, Windows)
The easiest way to run the entire stack without installing local dependencies.

1.  **Build and Run**:
    ```bash
    docker-compose up --build
    ```
2.  **Access**:
    The application will be available at `http://localhost:8080`.

> [!NOTE]
> The Docker container is standalone. If you wish to use different coupon files, you can swap the files in the `./requirements` folder and restart the container.

---

### Option 2: Local Development
If you prefer to run the components separately on your host machine.

#### Prerequisites
- [Go 1.25+](https://go.dev/)
- [Node.js & npm](https://nodejs.org/)

#### 1. Build the Frontend
```bash
cd web
npm install
npm run build
```
This will generate the `dist/` folder that the backend will serve.

#### 2. Run the Backend
From the root directory:
```bash
go run cmd/server/main.go [flags]
```

**Common Flags:**
- `-port`: The port the server will listen on (default: `8080`).
- `-coupon-files`: Comma-separated list of paths to your `.gz` coupon files (default: `requirements/couponbase1.gz,requirements/couponbase2.gz,requirements/couponbase3.gz`).

**Example for different OS paths:**
- **Windows**: `go run cmd/server/main.go -port 9000 -coupon-files "C:\coupons\file1.gz,C:\coupons\file2.gz"`
- **Linux/macOS**: `go run cmd/server/main.go -port 8080 -coupon-files "/tmp/coupons/1.gz,/tmp/coupons/2.gz"`

---

## 🔌 API Reference
The API documentation is based on OpenAPI 3.1.
- **UI Documentation**: [API Docs](https://orderfoodonline.deno.dev/public/openapi.html)
- **Spec File**: `api/openapi.yaml`

### Key Endpoints
- `GET /product`: List all products.
- `POST /order`: Place an order (**Header** `api_key: apitest` required).

#### Example Order Request
```json
{
  "items": [
    { "productId": "1", "quantity": 2 }
  ],
  "couponCode": "HAPPYHOURS"
}
```

## 🎨 Design Reference
- **Figma Design**: `design.fig`
- **Typography**: [Red Hat Text](https://fonts.google.com/specimen/Red+Hat+Text)
- **Colors**: HSL-based palette defined in `web/src/index.css`.

## 📈 Future Roadmap
- **Persistence**: Add PostgreSQL/MySQL support.
- **Caching**: Integrate Redis for rapid product lookups.
- **Async Processing**: Use Kafka or RabbitMQ for order fulfillment.
- **Observability**: Add Prometheus metrics and OpenTelemetry tracing.
