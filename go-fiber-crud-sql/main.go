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

	newProduct := &Product{
		Name:  "Iphone 14",
		Price: 13000,
		QuantityL 9
	}
}



func createProductDB(db *gorm.DB) {
	result := db.Create(product)

	if result.Error != nil {
		panic("failed to create product")
	}
	fmt.Println("Product created successfully")
}

func getProductAllDB(db *gorm.DB) {
    var product[]product

    result := db.Find(&product)

    if result.Error != nil {
        panic("failed to get products")
    }
}

func getProductByIDDB(db *gorm.DB, id int) {
    var product product

    result := db.First(&product, id)

    if result.Error != nil {
        panic("failed to get product")
    }

	fmt.Println("Get product by ID successful")
	fmt.Println(product)
}

func getProductByNameDB(db *gorm.DB, id int) {
    var product product

    result := db.Where("name = ?", name).First(&product)

    if result.Error != nil {
        panic("failed to get product")
    }

	fmt.Println("Get product by ID successful")
	fmt.Println(product)
}

func UpdateProductByIDDB(db *gorm.DB, id int) {
    result := db.Model(&product).Where("id = ?", id).Update(product)

    if result.Error != nil {
        panic("failed to get product")
    }

	fmt.Println("Get product by ID successful")
}

func DeleteProductByIDDB(db *gorm.DB, id int) {
    result := db.Delete(&product, id)

    if result.Error != nil {
        panic("failed to get product")
    }

	fmt.Println("Get product by ID successful")
}