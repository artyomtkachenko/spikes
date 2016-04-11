package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("failed to connect database")
	}
	/* defer db.Close() */

	// Create

	db.AutoMigrate(Product{})
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.Debug().First(&product, 1)           // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	//Print

	fmt.Printf("%d\n", product.Price)

	// Delete - delete product
	db.Delete(&product)
}
