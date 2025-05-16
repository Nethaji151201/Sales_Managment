# sales_report_master

## Tech stack

- Language : Go 2.3+
- CSV : File/Customer.CSV
- Schedular : Fix the time to run the data refresh go routines

## External Package

- http
- Toml
- UUID

### clone URL 

- git@github.com:Nethaji151201/Sales_Managment.git

### **1️⃣ Install Dependencies**

```sh
go mod tidy
```

## sample Request

```sh
curl http://localhost:29915/buyproduct?fromdate=2024-01-01&todate=2025-05-06

curl http://localhost:29915/revenueByCat?fromdate=2024-01-01&todate=2025-05-06

curl http://localhost:29915/revenueByRegion?fromdate=2024-01-01&todate=2025-05-06

curl http://localhost:29915/totalRevenue?fromdate=2024-01-01&todate=2025-05-06

curl http://localhost:29915/totalRevenuebyRange?fromdate=2024-01-01&todate=2025-05-06&rangetype=A


All API has GET Method

```

## sample Response 1
- Endpoint :  /buyproduct?fromdate=2024-01-01&todate=2025-05-06

```json
{"from_date":"2024-01-01","to_date":"2025-05-06","salesRepByProduct":[{"product_name":"shoee","totalRevenueWithDis":"2036.2","totalRevenueWithoutDis":"2145.5"}]}

```

### sample Response 2
- Endpoint :  /revenueByCat?fromdate=2024-01-01&todate=2025-05-06

```json
{"from_date":"2024-01-01","to_date":"2025-05-06","salesRepByProduct":[{"catagiryName":"shoee","totalRevenueWithDis":"2036.2","totalRevenueWithoutDis":"2145.5"}]}

```

### sample Response 3
- Endpoint :  /revenueByRegion?start_date=2016-01-01&end_date=2025-05-06

```json
{"from_date":"2024-01-01","to_date":"2025-05-06","salesRepByProduct":[{"regionName":"shoee","totalRevenueWithDis":"2036.2","totalRevenueWithoutDis":"2145.5"}]}

```

### sample Response 4
- Endpoint :  /totalRevenue?start_date=2016-01-01&end_date=2025-05-06

```json
{"from_date":"2024-01-01","to_date":"2025-05-06","totalRevenueWithDis":"4285.25","totalRevenueWithoutDis":"6547.25"}

```


### sample Response 5
- Endpoint :  /totalRevenuebyRange?fromdate=2024-01-01&todate=2025-05-06&rangetype=A

```json
{"from_date":"2024-01-01","to_date":"2025-05-06","yearWiseRevenue":[{"year":"2025","totalRevenueWithDis":"1524.25","totalRevenueWithoutDis":"1664.8"}],"monthWiseRevenue":[{"year":"2025","month":"JAN","totalRevenueWithDis":"1524.25","totalRevenueWithoutDis":"1664.8"}],"quaterWiseRevenue":[{"year":"2025","quater":"1","totalRevenueWithDis":"1524.25","totalRevenueWithoutDis":"1664.8"}]}

```


