package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// createCategoryDB เพิ่มประเภทสินค้า
func createCategoryDB(db *gorm.DB, name string, description string) (*Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category := &Category{
		Name:        name,
		Description: description,
	}

	if err := db.Create(category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// deleteCategoryByIDDB ลบประเภทสินค้า
// เงื่อนไข: ตรวจสอบความสัมพันธ์ก่อนลบ - ห้ามลบประเภทสินค้าที่ยังมีสินค้า
func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return fmt.Errorf("category ID %d not found: %w", id, err)
	}

	// ตรวจสอบสินค้าที่อยู่ในหมวดหมู่นี้
	var productCount int64
	if err := db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to check products for category: %w", err)
	}

	if productCount > 0 {
		return fmt.Errorf("ห้ามลบประเภทสินค้าที่ยังมีสินค้า: หมวดหมู่ ID %d มีสินค้าอยู่ %d รายการ", id, productCount)
	}

	if err := db.Delete(&category).Error; err != nil {
		return fmt.Errorf("failed to delete category ID %d: %w", id, err)
	}

	return nil
}
