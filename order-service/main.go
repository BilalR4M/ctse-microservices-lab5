package main

import (
	"database/context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Order struct {
	ID         int       `json:"id"`
	ItemID     int       `json:"item_id"`
	Quantity   int       `json:"quantity"`
	CustomerID string    `json:"customer_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func main() {
	time.Sleep(5 * time.Second) // wait for DB briefly

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "postgres"
		dbPort = "5432"
		dbUser = "postgres"
		dbPass = "password"
		dbName = "microservices_db"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var db *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		db, err = pgxpool.Connect(context.Background(), dsn)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to DB, retrying... (%v/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Unable to connect to pg DB: %v\n", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL DB.")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "order-service"})
	})

	r.GET("/orders", func(c *gin.Context) {
		rows, err := db.Query(context.Background(), "SELECT id, item_id, quantity, customer_id, status, created_at FROM order_schema.orders ORDER BY id DESC")
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(&o.ID, &o.ItemID, &o.Quantity, &o.CustomerID, &o.Status, &o.CreatedAt); err != nil {
				continue
			}
			orders = append(orders, o)
		}
		c.JSON(200, orders)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		var o Order
		err := db.QueryRow(context.Background(), "SELECT id, item_id, quantity, customer_id, status, created_at FROM order_schema.orders WHERE id = $1", id).Scan(&o.ID, &o.ItemID, &o.Quantity, &o.CustomerID, &o.Status, &o.CreatedAt)
		if err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(200, o)
	})

	r.POST("/orders", func(c *gin.Context) {
		var req struct {
			ItemID     int    `json:"item_id"`
			Quantity   int    `json:"quantity"`
			CustomerID string `json:"customer_id"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Validate Item exists (Step 8 inter-service call placeholder, we do it in Step 8. wait let's just do it directly so it works)
		itemSvcURL := os.Getenv("ITEM_SERVICE_URL")
		if itemSvcURL != "" {
			resp, err := http.Get(fmt.Sprintf("%s/items/%d", itemSvcURL, req.ItemID))
			if err != nil || resp.StatusCode != 200 {
				c.JSON(400, gin.H{"error": "Invalid Item ID"})
				return
			}
		}

		var newID int
		err := db.QueryRow(
			context.Background(),
			"INSERT INTO order_schema.orders (item_id, quantity, customer_id) VALUES ($1, $2, $3) RETURNING id",
			req.ItemID, req.Quantity, req.CustomerID,
		).Scan(&newID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		// fetch created
		var o Order
		db.QueryRow(context.Background(), "SELECT id, item_id, quantity, customer_id, status, created_at FROM order_schema.orders WHERE id = $1", newID).Scan(&o.ID, &o.ItemID, &o.Quantity, &o.CustomerID, &o.Status, &o.CreatedAt)

		c.JSON(201, o)
	})

	// PUT endpoint for Payment to update status
	r.PUT("/orders/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Status string `json:"status"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE order_schema.orders SET status = $1 WHERE id = $2", req.Status, id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(200, gin.H{"message": "Status updated"})
	})

	r.Run(":8082")
}
