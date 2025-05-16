package database

import (
	"customersales/apps/models"
	"customersales/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func ConnectDatabase() error {
	cf := config.GetConfig()
	fmt.Println(cf)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cf.Database.User,
		cf.Database.Pass,
		cf.Database.Host,
		cf.Database.Port,
		cf.Database.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	err = db.AutoMigrate(
		&models.Customer{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatalf("Initial table creation failed: %v", err)
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`
			ALTER TABLE orders 
			ADD CONSTRAINT fk_orders_customer 
			FOREIGN KEY (customer_id) 
			REFERENCES customers(customer_id)
		`).Error

		if err != nil {
			log.Println(err)
		}
		err = tx.Exec(`
			ALTER TABLE order_items 
			ADD CONSTRAINT fk_order_items_order 
			FOREIGN KEY (order_id) 
			REFERENCES orders(order_id)
		`).Error
		if err != nil {
			log.Println(err)
		}

		err = tx.Exec(`
			ALTER TABLE order_items 
			ADD CONSTRAINT fk_order_items_product 
			FOREIGN KEY (product_id) 
			REFERENCES products(product_id)
		`).Error
		if err != nil {
			log.Println(err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to add foreign key constraints: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
		return err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	GDB = db
	log.Println("Database setup completed successfully")
	return nil
}
