# SE4010 Microservices Lab 5

This repository contains the implementation of the SE4010 Microservices Lab System. The project follows a strictly Docker-only development methodology, integrating multiple polyglot microservices, a central PostgreSQL database, and an API Gateway.

## Architecture & Technology Stack
The lab uses **Option B** for its technology stack:

| Component       | Technology               | Port   | Description                                                                              |
|-----------------|--------------------------|--------|------------------------------------------------------------------------------------------|
| **Item Service**    | Node.js / Express        | `8081` | Manages store items. Exposes endpoints to retrieve and create items.                     |
| **Order Service**   | Go / Gin                   | `8082` | Manages customer orders. Interacts with the Item Service to validate item availability.  |
| **Payment Service** | Python / Flask             | `8083` | Handles payments. Interacts with the Order Service to fetch and update order status.       |
| **API Gateway** | Kong (DB-less mode)      | `8080` | A central reverse proxy routing all traffic from clients to the underlying microservices. |
| **Database**    | PostgreSQL               | `5432` | A single container sharing three distinct logical schemas (`item_schema`, `order_schema`, `payment_schema`). |

All microservices are isolated inside a shared Docker network (`microservices-net`) and communicate with each other exclusively via absolute Docker service hostnames (e.g., `http://item-service:8081`).

## Prerequisites
- Docker
- Docker Compose v2

*Note: No local installation of Node.js, Go, Python, or PostgreSQL is required. Everything runs natively inside Docker containers.*

## Running the Application

1. **Clone the repository:**
   ```bash
   git clone https://github.com/BilalR4M/ctse-microservices-lab5.git
   cd ctse-microservices-lab5
   ```

2. **Build and start the services:**
   ```bash
   docker-compose up -d --build
   ```

3. **Verify running containers:**
   ```bash
   docker-compose ps
   ```
   *You should see `postgres`, `kong`, `item-service`, `order-service`, and `payment-service` all running.*

## API Endpoints

All requests should be sent to the API Gateway at `http://localhost:8080`.

### Item Service (`/items`)
- `GET /items`: Fetch all items.
- `GET /items/:id`: Fetch item by ID.
- `POST /items`: Create a new item.
  ```json
  { "name": "Laptop" }
  ```

### Order Service (`/orders`)
- `GET /orders`: Fetch all orders.
- `GET /orders/:id`: Fetch order by ID.
- `POST /orders`: Create a new order *(validates item exists via Item Service)*.
  ```json
  { "item_id": 1, "quantity": 2, "customer_id": "cust-001" }
  ```
- `PUT /orders/:id/status`: Update order status.
  ```json
  { "status": "SHIPPED" }
  ```

### Payment Service (`/payments`)
- `GET /payments`: Fetch all payments.
- `GET /payments/:id`: Fetch payment by ID.
- `POST /payments/process`: Process payment for an order *(updates order status in Order Service to `PAID` via API call)*.
  ```json
  { "order_id": 1, "amount": 1999.98 }
  ```

## Working with Postman
A Postman collection is included in the root directory: `postman_collection.json`. 
Import this file into your Postman application to instantly access and test all pre-configured API endpoints.

## Troubleshooting & Logs
To debug specific services or view terminal output, use the Docker Compose logging feature:

```bash
docker-compose logs -f item-service
docker-compose logs -f order-service
docker-compose logs -f payment-service
```
