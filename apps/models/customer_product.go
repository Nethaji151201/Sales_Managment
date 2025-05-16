package models

import "time"

type Customer struct {
	CustomerID string `gorm:"primaryKey;type:varchar(50)"`
	Name       string
	Email      string
	Address    string
}

type Product struct {
	ProductID string `gorm:"primaryKey;type:varchar(50)"`
	Name      string
	Category  string
}

type Order struct {
	OrderID      string `gorm:"primaryKey;type:varchar(50)"`
	CustomerID   string `gorm:"type:varchar(50);not null"`
	Region       string
	DateOfSale   time.Time
	PaymentType  string
	ShippingCost float64
}

type OrderItem struct {
	ID           uint   `gorm:"primaryKey"`
	OrderID      string `gorm:"type:varchar(50);not null"`
	ProductID    string `gorm:"type:varchar(50);not null"`
	QuantitySold int
	UnitPrice    float64
	Discount     float64
}

type RevenueResult struct {
	ProductName            string `json:"product_name,omitempty" gorm:"ProductName"`
	CatagiryName           string `json:"catagiryName,omitempty" gorm:"CatagiryName"`
	RegionName             string `json:"regionName,omitempty" gorm:"RegionName"`
	TotalRevenueWithDis    string `json:"totalRevenueWithDis" gorm:"TotalRevenueWithDis"`
	TotalRevenueWithoutDis string `json:"totalRevenueWithoutDis" gorm:"TotalRevenueWithoutDis"`
}

type RevenueByRangeStruct struct {
	Month                  string `json:"month,omitempty" gorm:"Month"`
	Year                   string `json:"year,omitempty" gorm:"Year"`
	Quater                 string `json:"quater,omitempty" gorm:"Quater"`
	TotalRevenueWithDis    string `json:"totalRevenueWithDis,omitempty" gorm:"TotalRevenueWithDis"`
	TotalRevenueWithoutDis string `json:"totalRevenueWithoutDis,omitempty" gorm:"TotalRevenueWithoutDis"`
}

type ProductStruct struct {
	FromDate               string                 `json:"from_date"`
	ToDate                 string                 `json:"to_date"`
	SalesRepByProduct      []RevenueResult        `json:"salesRepByProduct,omitempty" `
	SalesRepByCatagiry     []RevenueResult        `json:"salesRepByCatagiry,omitempty" `
	SalesRepByRegion       []RevenueResult        `json:"salesRepByRegion,omitempty" `
	YearWiseRevenue        []RevenueByRangeStruct `json:"yearWiseRevenue,omitempty" `
	MonthWiseRevenue       []RevenueByRangeStruct `json:"monthWiseRevenue,omitempty" `
	QuaterWiseRevenue      []RevenueByRangeStruct `json:"quaterWiseRevenue,omitempty" `
	TotalRevenueWithDis    string                 `json:"totalRevenueWithDis,omitempty" gorm:"TotalRevenueWithDis"`
	TotalRevenueWithoutDis string                 `json:"totalRevenueWithoutDis,omitempty" gorm:"TotalRevenueWithoutDis"`
	Status                 string                 `json:"status"`
	ErrMsg                 string                 `json:"errMsg"`
}
