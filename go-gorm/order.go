package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Allowed order statuses
const (
	StatusPending   = "Pending"
	StatusCompleted = "Completed"
	StatusCancelled = "Cancelled"
)

func validateStatus(status string) (string, error) {
	cleanStatus := strings.TrimSpace(status)
	switch strings.ToLower(cleanStatus) {
	case "pending":
		return StatusPending, nil
	case "completed":
		return StatusCompleted, nil
	case "cancelled", "canceled":
		return StatusCancelled, nil
	default:
		return "", fmt.Errorf("สถานะต้องเป็น Pending, Completed หรือ Cancelled เท่านั้น (ได้รับ: '%s')", status)
	}
}

// createOrderDB เพิ่มคำสั่งซื้อ
// ต้องตรวจสอบลูกค้าและสินค้า ตรวจสอบสต็อก คำนวณ TotalPrice = Price * Quantity
// ลดสต็อก และใช้ Transaction ให้การเพิ่มคำสั่งซื้อกับการลดสต็อกสำเร็จหรือยกเลิกพร้อมกัน
func createOrderDB(db *gorm.DB, customerID uint, productID uint, quantity int) (*Order, error) {
	if quantity <= 0 {
		return nil, errors.New("จำนวนสินค้าต้องมากกว่า 0")
	}

	var newOrder Order

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. ตรวจสอบข้อมูลลูกค้า
		var customer Customer
		if err := tx.First(&customer, customerID).Error; err != nil {
			return fmt.Errorf("ไม่พบข้อมูลลูกค้า ID %d: %w", customerID, err)
		}

		// 2. ตรวจสอบข้อมูลสินค้า
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return fmt.Errorf("ไม่พบข้อมูลสินค้า ID %d: %w", productID, err)
		}

		// 3. ตรวจสอบสต็อก
		if product.Quantity < quantity {
			return fmt.Errorf("สินค้า '%s' มีสต็อกไม่เพียงพอ (คงเหลือ %d ชิ้น, ต้องการสั่งซื้อ %d ชิ้น)",
				product.Name, product.Quantity, quantity)
		}

		// 4. คำนวณ TotalPrice = Price * Quantity
		totalPrice := product.Price * float64(quantity)

		// 5. ลดสต็อกสินค้า
		product.Quantity -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to deduct stock: %w", err)
		}

		// 6. บันทึกคำสั่งซื้อ
		newOrder = Order{
			CustomerID: customerID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     StatusPending,
		}

		if err := tx.Create(&newOrder).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Preload Customer and Product
	db.Preload("Customer").Preload("Product").Preload("Product.Category").First(&newOrder, newOrder.ID)
	return &newOrder, nil
}

// getOrdersByCustomerDB ดึงคำสั่งซื้อของลูกค้า
func getOrdersByCustomerDB(db *gorm.DB, customerID uint) ([]Order, error) {
	var orders []Order
	if err := db.Preload("Customer").Preload("Product").Where("customer_id = ?", customerID).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders for customer %d: %w", customerID, err)
	}
	return orders, nil
}

// getOrdersByStatusDB ดึงคำสั่งซื้อตามสถานะ (Pending / Completed / Cancelled)
func getOrdersByStatusDB(db *gorm.DB, status string) ([]Order, error) {
	validStatus, err := validateStatus(status)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := db.Preload("Customer").Preload("Product").Where("status = ?", validStatus).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders with status '%s': %w", validStatus, err)
	}
	return orders, nil
}

// getOrdersByDateRangeDB ดึงคำสั่งซื้อตามช่วงวันที่
func getOrdersByDateRangeDB(db *gorm.DB, startDate time.Time, endDate time.Time) ([]Order, error) {
	var orders []Order
	if err := db.Preload("Customer").Preload("Product").Where("created_at >= ? AND created_at <= ?", startDate, endDate).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders by date range [%s to %s]: %w",
			startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), err)
	}
	return orders, nil
}

// getOrdersWithDetailsDB ดึงคำสั่งซื้อพร้อมข้อมูลลูกค้าและสินค้า
func getOrdersWithDetailsDB(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Preload("Customer").Preload("Product").Preload("Product.Category").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders with details: %w", err)
	}
	return orders, nil
}

// updateOrderStatusDB เปลี่ยนสถานะคำสั่งซื้อ
// เงื่อนไข: สถานะต้องเป็น Pending, Completed หรือ Cancelled เท่านั้น
func updateOrderStatusDB(db *gorm.DB, id uint, status string) (*Order, error) {
	validStatus, err := validateStatus(status)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return nil, fmt.Errorf("order ID %d not found: %w", id, err)
	}

	order.Status = validStatus
	if err := db.Save(&order).Error; err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	db.Preload("Customer").Preload("Product").First(&order, order.ID)
	return &order, nil
}
