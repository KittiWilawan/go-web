package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// createCustomerDB เพิ่มลูกค้า
func createCustomerDB(db *gorm.DB, name string, email string, phone string) (*Customer, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" {
		return nil, errors.New("ชื่อลูกค้าต้องไม่เป็นค่าว่าง (customer name cannot be empty)")
	}
	if email == "" {
		return nil, errors.New("อีเมลต้องไม่เป็นค่าว่าง (customer email cannot be empty)")
	}

	customer := &Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	if err := db.Create(customer).Error; err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// getCustomerByEmailDB ค้นหาลูกค้าด้วย Email
func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	email = strings.TrimSpace(email)
	var customer Customer
	if err := db.Where("email = ?", email).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer with email '%s' not found: %w", email, err)
	}
	return &customer, nil
}

// updateCustomerByIDDB แก้ไขลูกค้า
func updateCustomerByIDDB(db *gorm.DB, id uint, name string, email string, phone string) (*Customer, error) {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return nil, fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" {
		return nil, errors.New("customer name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("customer email cannot be empty")
	}

	customer.Name = name
	customer.Email = email
	customer.Phone = phone

	if err := db.Save(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to update customer ID %d: %w", id, err)
	}

	return &customer, nil
}

// deleteCustomerByIDDB ลบลูกค้าตาม ID
// เงื่อนไข: ตรวจสอบความสัมพันธ์ก่อนลบ - ห้ามลบลูกค้าที่มีประวัติคำสั่งซื้อ
func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return fmt.Errorf("customer ID %d not found: %w", id, err)
	}

	// ตรวจสอบประวัติคำสั่งซื้อ
	var orderCount int64
	if err := db.Model(&Order{}).Where("customer_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check customer orders: %w", err)
	}

	if orderCount > 0 {
		return fmt.Errorf("ห้ามลบลูกค้าที่มีประวัติคำสั่งซื้อ: ลูกค้า ID %d มีคำสั่งซื้ออยู่ %d รายการ", id, orderCount)
	}

	if err := db.Delete(&customer).Error; err != nil {
		return fmt.Errorf("failed to delete customer ID %d: %w", id, err)
	}

	return nil
}
