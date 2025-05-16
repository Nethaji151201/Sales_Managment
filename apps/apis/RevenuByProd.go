package apis

import (
	"customersales/apps/models"
	"customersales/apps/utils"
	database "customersales/migrations"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Sample Request json

/* {"from_date" : "2024-01-01", "to_date":"2024-04-30" }*/

/* {"from_date" : "2024-01-01", "to_date":"2024-04-30" }*/

func BuyProductDetails(w http.ResponseWriter, r *http.Request) {
	log := new(utils.Logger)
	log.SetSid(r)
	log.Log(utils.INFO, "BuyProductDetails (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	var lErr error
	var lRespRec models.ProductStruct

	switch r.Method {
	case http.MethodGet:
		{
			lRespRec.Status = "S"
			lRespRec.FromDate = r.URL.Query().Get("fromdate")
			lRespRec.ToDate = r.URL.Query().Get("todate")
			if lRespRec.FromDate == "" || lRespRec.ToDate == "" {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "Date Range mismatch"
				log.Log(utils.ERROR, "Date Range mismatch")
				goto marshal
			}

			sDate, lErr := time.Parse("2006-01-02", lRespRec.FromDate)
			if lErr != nil {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "RevenuByCat - 003" + lErr.Error()
				log.Log(utils.ERROR, "RevenuByCat - 003"+lErr.Error())
				goto marshal
			}

			eDate, err := time.Parse("2006-01-02", lRespRec.ToDate)
			if err != nil {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "RevenuByCat - 004" + lErr.Error()
				log.Log(utils.ERROR, "RevenuByCat - 004"+lErr.Error())
				goto marshal
			}

			if sDate.After(eDate) {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "ToDate should be greater than FromDate"
				log.Log(utils.ERROR, "ToDate should be greater than FromDate")
				goto marshal
			}

			lRespRec.SalesRepByProduct, lErr = BuyProduct(log, lRespRec)
			if lErr != nil {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "BuyProductDetails-003 " + lErr.Error()
				log.Log(utils.ERROR, "BuyProductDetails-003 ", lErr.Error())

			}

		}
	default:
		{
			lRespRec.Status = "E"
			lRespRec.ErrMsg = "Invalid Method"
		}
	}
marshal:
	lData, lErr := json.Marshal(lRespRec)
	if lErr != nil {
		fmt.Fprintf(w, "Error taking data"+lErr.Error())
	} else {
		fmt.Fprintf(w, string(lData))
	}
	log.Log(utils.INFO, "BuyProductDetails (-)")
}

func BuyProduct(log *utils.Logger, lRespRec models.ProductStruct) ([]models.RevenueResult, error) {
	log.Log(utils.INFO, "BuyProduct (+)")

	var results []models.RevenueResult

	lErr := database.GDB.Table("order_items oi").
		Select(`p.name as ProductName, sum(quantity_sold * unit_price * (1 - discount)) as TotalRevenueWithDis,sum(quantity_sold * unit_price) as TotalRevenueWithoutDis`).
		Joins(`join products p 
				on oi.product_id = p.product_id `).
		Joins(`join orders o 
				on o.order_id = oi.order_id `).
		Where(`o.date_of_sale between ? and ?`, lRespRec.FromDate, lRespRec.ToDate).
		Group(`p.product_id`).
		Scan(&results).Error

	if lErr != nil {
		log.Log(utils.ERROR, "GR005", lErr.Error())
		return results, fmt.Errorf("failed to calculate total revenue: %w", lErr.Error())
	}

	log.Log(utils.INFO, "BuyProduct (-)")
	return results, nil
}

func DataRefersh(log *utils.Logger, path string) error {
	log.Log(utils.INFO, "DataRefersh (+)")

	lFile, err := os.Open(path)
	if err != nil {
		log.Log(utils.ERROR, "DR001", err.Error())
		return err
	}
	defer lFile.Close()

	reader := csv.NewReader(lFile)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		log.Log(utils.ERROR, "DR002", err.Error())
		return err
	}

	for _, row := range rows[1:] {
		orderID := row[0]
		productID := row[1]
		customerID := row[2]

		quantity, err := strconv.Atoi(row[7])
		if err != nil {
			log.Log(utils.ERROR, "DR003", err.Error())
			return err
		}

		unitPrice, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR004", err.Error())
			return err
		}

		discount, err := strconv.ParseFloat(row[9], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR005", err.Error())
			return err
		}

		shippingCost, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR006", err.Error())
			return err
		}

		dateOfSale, err := time.Parse("2006-01-02", row[6])
		if err != nil {
			log.Log(utils.ERROR, "DR007", err.Error())
			return err
		}

		customer := models.Customer{CustomerID: customerID}
		err = database.GDB.Where("customer_id = ?", customerID).FirstOrCreate(&customer, models.Customer{
			Name:    row[12],
			Email:   row[13],
			Address: row[14],
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR008", err.Error())
			return err
		}

		product := models.Product{ProductID: productID}
		err = database.GDB.Where("product_id = ?", productID).FirstOrCreate(&product, models.Product{
			Name:     row[3],
			Category: row[4],
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR009", err.Error())
			return err
		}

		order := models.Order{OrderID: orderID}
		err = database.GDB.Where("order_id = ?", orderID).FirstOrCreate(&order, models.Order{
			CustomerID:   customerID,
			Region:       row[5],
			DateOfSale:   dateOfSale,
			PaymentType:  row[11],
			ShippingCost: shippingCost,
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR010", err.Error())
			return err
		}

		orderItem := models.OrderItem{
			OrderID:      orderID,
			ProductID:    productID,
			QuantitySold: quantity,
			UnitPrice:    unitPrice,
			Discount:     discount,
		}

		err = database.GDB.Create(&orderItem).Error
		if err != nil {
			log.Log(utils.ERROR, "DR011", err.Error())
			return err
		}
	}

	log.Log(utils.INFO, "CSV data processed and refreshed successfully")
	log.Log(utils.INFO, "DataRefersh (-)")
	return nil
}
