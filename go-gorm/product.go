package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createProductDB(db *gorm.DB, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("ราคาสินค้าต้องไม่ติดลบ (price cannot be negative)")
	}
	if quantity < 0 {
		return nil, errors.New("จำนวนสต็อกต้องไม่ติดลบ (quantity cannot be negative)")
	}

	// ตรวจสอบว่ามี Category อยู่จริง
	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("ไม่พบประเภทสินค้า ID %d: %w", categoryID, err)
	}

	product := &Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	}

	if err := db.Create(product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	db.Preload("Category").First(product, product.ID)
	return product, nil
}

// getAllProductsDB ค้นหาสินค้าทั้งหมด
func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	return products, nil
}

// getProductByIDDB ค้นหาสินค้าจาก ID
func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}
	return &product, nil
}

// getProductByNameDB ค้นหาชื่อบางส่วนด้วย LIKE
func getProductByNameDB(db *gorm.DB, keyword string) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by name keyword '%s': %w", keyword, err)
	}
	return products, nil
}

// getProductsByCategoryDB ค้นหาสินค้าตาม CategoryID
func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products by category %d: %w", categoryID, err)
	}
	return products, nil
}

// getProductsByPriceRangeDB ค้นหาสินค้าตามช่วงราคา Min–Max
func getProductsByPriceRangeDB(db *gorm.DB, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("price >= ? AND price <= ?", minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products in price range [%.2f, %.2f]: %w", minPrice, maxPrice, err)
	}
	return products, nil
}

// getProductsLowStockDB ค้นหาสินค้าที่สต็อกน้อยกว่าค่ากำหนด
func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("quantity < ?", threshold).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch low stock products (< %d): %w", threshold, err)
	}
	return products, nil
}

// getProductsByCategoryAndPriceDB ค้นหาตามประเภทและช่วงราคา
func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ? AND price >= ? AND price <= ?", categoryID, minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products by category %d and price range [%.2f, %.2f]: %w", categoryID, minPrice, maxPrice, err)
	}
	return products, nil
}

// getProductsSortedByPriceDB เรียงราคา ASC/DESC
func getProductsSortedByPriceDB(db *gorm.DB, sortOrder string) ([]Product, error) {
	upperOrder := strings.ToUpper(strings.TrimSpace(sortOrder))
	if upperOrder != "ASC" && upperOrder != "DESC" {
		upperOrder = "ASC"
	}

	var products []Product
	if err := db.Preload("Category").Order("price " + upperOrder).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sorted products: %w", err)
	}
	return products, nil
}

// getTopExpensiveProductsDB สินค้าแพงที่สุด N รายการ
func getTopExpensiveProductsDB(db *gorm.DB, n int) ([]Product, error) {
	if n <= 0 {
		n = 3
	}
	var products []Product
	if err := db.Preload("Category").Order("price DESC").Limit(n).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch top %d expensive products: %w", n, err)
	}
	return products, nil
}

// updateProductByIDDB แก้ไขสินค้าตาม ID
func updateProductByIDDB(db *gorm.DB, id uint, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("ราคาสินค้าต้องไม่ติดลบ")
	}
	if quantity < 0 {
		return nil, errors.New("สต็อกสินค้าต้องไม่ติดลบ")
	}

	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("category ID %d not found: %w", categoryID, err)
	}

	product.Name = name
	product.Price = price
	product.Quantity = quantity
	product.CategoryID = categoryID

	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product ID %d: %w", id, err)
	}

	db.Preload("Category").First(&product, product.ID)
	return &product, nil
}

// updateProductPriceDB แก้ไขราคา
func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) (*Product, error) {
	if newPrice < 0 {
		return nil, errors.New("ราคาต้องไม่ติดลบ (price cannot be negative)")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	product.Price = newPrice
	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product price: %w", err)
	}

	return &product, nil
}

// updateProductStockDB เพิ่ม/ลดสต็อก (ส่งค่า quantityChange เป็นบวกเพื่อเพิ่ม หรือลบเพื่อลด)
func updateProductStockDB(db *gorm.DB, id uint, quantityChange int) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	newStock := product.Quantity + quantityChange
	if newStock < 0 {
		return nil, fmt.Errorf("สต็อกสินค้าต้องไม่ติดลบ (สต็อกเดิม: %d, เปลี่ยนแปลง: %d, สต็อกใหม่: %d)", product.Quantity, quantityChange, newStock)
	}

	product.Quantity = newStock
	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update stock for product ID %d: %w", id, err)
	}

	return &product, nil
}

// softDeleteProductDB ทำ Soft Delete สินค้า
func softDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if err := db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to soft delete product ID %d: %w", id, err)
	}
	return nil
}

// hardDeleteProductDB ลบสินค้าถาวร
// กติกาส่วนที่ 5: สินค้าที่อยู่ในประวัติคำสั่งซื้อ ห้ามลบถาวร เพื่อรักษาความสมบูรณ์ของประวัติคำสั่งซื้อ
func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	// ตรวจสอบประวัติคำสั่งซื้อ
	var orderCount int64
	if err := db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check orders for product ID %d: %w", id, err)
	}

	if orderCount > 0 {
		return fmt.Errorf("ห้ามลบสินค้าถาวร: สินค้า ID %d มีประวัติคำสั่งซื้ออยู่ในระบบ %d รายการ", id, orderCount)
	}

	if err := db.Unscoped().Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to hard delete product ID %d: %w", id, err)
	}
	return nil
}

// getDeletedProductsDB ดึงสินค้าที่ถูก Soft Delete
func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get soft deleted products: %w", err)
	}
	return products, nil
}

// restoreProductDB กู้คืนสินค้า
func restoreProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&product).Error; err != nil {
		return fmt.Errorf("soft deleted product ID %d not found: %w", id, err)
	}

	if err := db.Unscoped().Model(&product).Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("failed to restore product ID %d: %w", id, err)
	}
	return nil
}
