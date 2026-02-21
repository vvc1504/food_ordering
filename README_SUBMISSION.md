# Food Ordering API Server

A robust, extensible Food Ordering API built in Go, adhering to the provided OpenAPI 3.1 specification.

## 🚀 Key Features

- **Sequential Coupon Validation**: High-performance engine that processes large `.gz` files with minimal memory footprint.
- **Production-Ready Architecture**: Layered design separating interfaces (spec) from implementation (impl).
- **Scalability**: Designed for stateless horizontal scaling.
- **Dockerized**: Ready for containerized deployment.
- **Strict OpenAPI Adherence**: Implements the mobile food ordering spec.

## 🏗️ Architecture Design (DDD Refactored)

The project follows a high-end DDD Aggregate pattern to ensure production-grade extensibility and scalability:

- **Aggregate Roots**: `ProductManager` and `OrderManager` act as entry points for their respective domains.
- **Async Pattern**: All domain logic is non-blocking, returning results and status signals via async channels (`<-chan`).
- **Granular Status**: Uses aggregate-specific status codes and types for precise error reporting, modeled after your design system.
- **Spec/Impl Pattern**: Domain interfaces are defined in `pkg/spec` (Aggregates) and concrete logic resides in `pkg/spec/impl`.
- **Aggregate Initialization**: Supports bulk initialization of the product aggregate via a dedicated `ProductInit` command.

## 🎟️ Coupon Validation Engine

The assignment requires validating coupons found in at least two of three large `.gz` files (totaling ~2.1GB gzipped).

### My Approach: Sequential Streaming Indexing
To handle this at scale without exhausting disk or memory:
1. **Streaming Decompression**: Files are never unzipped to disk. Instead, we use `compress/gzip` to stream data directly into memory.
2. **Sequential Processing**: Each file is processed one by one.
3. **Optimized Bitmask Indexing**:
   - Only strings of 8-10 characters are tracked.
   - We use intermediate tracking maps (`seenInFile1`, `seenInFile2`) to build a final `ValidCoupons` map.
   - Intermediate maps are garbage collected after initialization, leaving only the validated coupon set in memory.

## 🛠️ Setup & Running

### Prerequisites
- Go 1.21+ (if running locally)
- Docker & Docker Compose

### Running with Docker (Recommended)
```bash
docker-compose up --build
```
The server will be available at `http://localhost:8080`.

### Running Locally
```bash
go run cmd/server/main.go -port 8080 -coupon-dir ./requirements
```

## 🔌 API Endpoints

- `GET /product`: List all products.
- `GET /product/{id}`: Get product details.
- `POST /order`: Place an order (**Requires header** `api_key: apitest`).

### Example Order Request
```json
{
  "items": [
    { "productId": "1", "quantity": 2 },
    { "productId": "2", "quantity": 1 }
  ],
  "couponCode": "HAPPYHRS"
}
```

## 📈 Future Improvements
- **Persistent Database**: Replace in-memory repository with Postgres/MongoDB.
- **External Caching**: Use Redis for frequently accessed product data.
- **Async Order Processing**: Implement a message queue (RabbitMQ/Kafka) for order fulfillment.
- **Metrics & Tracing**: Add Prometheus metrics and OpenTelemetry tracing.
