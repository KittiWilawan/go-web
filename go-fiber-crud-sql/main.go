package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

const (
	host	 = "localhost"
	port	 = "3306"
	user	 = "root"
	password = ""
	dbname	 = "stock_db"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate the schema
	db.AutoMigrate(&Product{})

	fmt.Println(db)
	fmt.Println("Database connection successful")

	db.AutoMigrate(&Product{})
	fmt.Println("Database migration completed")
	fmt.Println("Database connection successful")
}