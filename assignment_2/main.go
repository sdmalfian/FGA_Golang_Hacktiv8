package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	ID           uint       `gorm:"primaryKey"`
	CustomerName string     `json:"customerName"`
	OrderedAt    time.Time  `json:"orderedAt"`
	Items        []Item     `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-"` // Gorm 1.20: Gunakan tipe gorm.DeletedAt untuk soft delete
}

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}

var db *gorm.DB

func main() {
	// db connection
	dsn := "host=localhost user=postgres password=postgres dbname=db_latihan port=5434 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// model migration
	err = db.AutoMigrate(&Order{}, &Item{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	r := gin.Default()

	// routes endpoint
	r.POST("/orders", createOrder)
	r.GET("/orders", getOrders)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)

	// start server at port 8080
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func getOrders(c *gin.Context) {
	var orders []Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	var response []gin.H
	for _, order := range orders {
		itemDetails := make([]gin.H, len(order.Items))
		for i, item := range order.Items {
			itemDetails[i] = gin.H{
				"itemCode":    item.Code,
				"description": item.Description,
				"quantity":    item.Quantity,
			}
		}

		orderDetails := gin.H{
			"id":           order.ID,
			"orderedAt":    order.OrderedAt,
			"customerName": order.CustomerName,
			"items":        itemDetails,
		}

		response = append(response, orderDetails)
	}

	c.JSON(http.StatusOK, response)
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var updatedOrder Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	order.CustomerName = updatedOrder.CustomerName
	order.OrderedAt = updatedOrder.OrderedAt

	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
