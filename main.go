package main

import (
	"customersales/apps/apis"
	"customersales/apps/utils"
	"customersales/config"
	database "customersales/migrations"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.OpenFile("./log/log"+time.Now().Format("02012006.15.04.05.000000000")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error open file ", err.Error())
	}
	defer file.Close()
	log.SetOutput(file)
	err = database.ConnectDatabase()
	if err != nil {
		log.Fatal("Error Occur in DB Connection : ", err)
	}
	go Schedular()

	lPort := config.GetConfig().Service.Port
	fmt.Println("Server Started : " + strconv.Itoa(lPort))

	//============================= API ==================================================
	http.HandleFunc("/buyproduct", apis.BuyProductDetails)
	http.HandleFunc("/revenueByCat", apis.RevenuByCat)
	http.HandleFunc("/revenueByRegion", apis.RevenuByRegion)
	http.HandleFunc("/totalRevenue", apis.TotalRevenue)
	http.HandleFunc("/totalRevenuebyRange", apis.RavanueByRange)

	srv := http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  30 * time.Second,
		Addr:         ":" + strconv.Itoa(lPort),
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
func Schedular() {
	log := new(utils.Logger)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Log(utils.INFO, "Data refersh ")
			err := apis.DataRefersh(log, "File/Customer.csv")
			if err != nil {
				log.Log(utils.ERROR, "Error Occur in schedular :", err.Error())
			}
		}
	}
}

func init() {
	config.LoadGlobalConfig("toml/config.toml")
}
