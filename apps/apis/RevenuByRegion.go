package apis

import (
	"customersales/apps/models"
	"customersales/apps/utils"
	database "customersales/migrations"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Sample Request json

/* {"from_date" : "2024-01-01", "to_date":"2024-04-30" }*/

/* {"from_date" : "2024-01-01", "to_date":"2024-04-30" }*/

func RevenuByRegion(w http.ResponseWriter, r *http.Request) {
	log := new(utils.Logger)
	log.SetSid(r)
	log.Log(utils.INFO, "RevenuByRegion (+)")

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
				lRespRec.ErrMsg = "RevenuByRegion - 003" + lErr.Error()
				log.Log(utils.ERROR, "RevenuByRegion - 003"+lErr.Error())
				goto marshal
			}

			eDate, err := time.Parse("2006-01-02", lRespRec.ToDate)
			if err != nil {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "RevenuByRegion - 004" + lErr.Error()
				log.Log(utils.ERROR, "RevenuByRegion - 004"+lErr.Error())
				goto marshal
			}

			if sDate.After(eDate) {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "ToDate should be greater than FromDate"
				log.Log(utils.ERROR, "ToDate should be greater than FromDate")
				goto marshal
			}

			lRespRec.SalesRepByCatagiry, lErr = RegionDetails(log, lRespRec)
			if lErr != nil {
				lRespRec.Status = "E"
				lRespRec.ErrMsg = "RevenuByRegion-003 " + lErr.Error()
				log.Log(utils.ERROR, "RevenuByRegion-003 ", lErr.Error())
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
	log.Log(utils.INFO, "RevenuByRegion (-)")
}

func RegionDetails(log *utils.Logger, lRespRec models.ProductStruct) ([]models.RevenueResult, error) {
	log.Log(utils.INFO, "RegionDetails (+)")

	var results []models.RevenueResult

	lErr := database.GDB.Table("order_items oi").
		Select("o.region as RegionName, sum(quantity_sold * unit_price * (1 - discount)) as TotalRevenueWithDis,sum(quantity_sold * unit_price) as TotalRevenueWithoutDis").
		Joins(`join products p 
				on oi.product_id = p.product_id `).
		Joins(`join orders o 
				on o.order_id = oi.order_id `).
		Where("o.date_of_sale between ? and ?", lRespRec.FromDate, lRespRec.ToDate).
		Group("o.region").
		Scan(&results).Error

	if lErr != nil {
		log.Log(utils.ERROR, "RegionDetails - 001", lErr.Error())
		return results, fmt.Errorf("failed to calculate total revenue: %w", lErr.Error())
	}

	log.Log(utils.INFO, "RegionDetails (-)")
	return results, nil
}
