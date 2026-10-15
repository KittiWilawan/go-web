package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Category model represents product categories
type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// Product model represents items in store
type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);not null" json:"name"`
	Price      float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity   int      `gorm:"not null;default:0" json:"quantity"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Orders     []Order  `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

// Customer model represents store customers
type Customer struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(100);not null" json:"name"`
	Email  string  `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone  string  `gorm:"type:varchar(20)" json:"phone"`
	Orders []Order `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

// Order model represents customer purchase orders
type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null" json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductID  uint     `gorm:"not null" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	TotalPrice float64  `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string   `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"` // Pending, Completed, Cancelled
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// InitDB connects to MySQL, creates store_db if not exists, and performs AutoMigrate
func InitDB() (*gorm.DB, error) {
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASS", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "store_db")

	// 1. Connect without db to ensure database exists
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)
	initDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbname)
	if err := initDB.Exec(createDBSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create database %s: %w", dbname, err)
	}

	sqlDB, err := initDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// 2. Connect to the store_db
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbname, err)
	}

	// 3. Auto Migrate all 4 models
	if err := db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{}); err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Database connection and AutoMigrate completed successfully.")
	return db, nil
}
