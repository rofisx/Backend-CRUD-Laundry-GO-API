package transaction

import (
	"challenge-goapi/config"
	"challenge-goapi/utils/getdate"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

func CreateTransaksiLaundry(c *gin.Context) {
	var newTrs TrsHeader
	var newResult Result
	err := c.ShouldBind(&newTrs)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkDate(newTrs.BillDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Tanggal Bill Harus YYYY-DD-MM HH24:MI:SS"})
		return
	}

	if !checkDate(newTrs.EntryDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Tanggal Entry Harus YYYY-DD-MM HH24:MI:SS"})
		return
	}

	if !checkDate(newTrs.FinishDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Tanggal Finish Harus YYYY-DD-MM HH24:MI:SS"})
		return
	}

	if !checkEmployeeIdExist(newTrs.EmployeeId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee Id Tidak Ada"})
		return
	}
	if !checkCustomerIdExist(newTrs.CustomerId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer Id Tidak Ada"})
		return
	}

	queryInsert := "INSERT INTO trs_laundry (id,billdate,entrydate,finishdate,employeeid,customerid) VALUES ($1,$2,$3,$4,$5,$6)"

	queryInsertDetail := "INSERT INTO trs_laundry_detail (billid,productid,qty) VALUES ($1,$2,$3) RETURNING id,billid,productid"

	selectPrice := `SELECT aa.price FROM mst_product as aa
					INNER JOIN
					trs_laundry_detail as bb
					ON aa.id = bb.productId
					WHERE bb.id = $1 AND bb.productid = $2 `

	newTrs.Id = createTrsId()
	fmt.Println(newTrs.Id)

	tx, errx := db.Begin()
	if errx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	_, errx = tx.Exec(queryInsert, &newTrs.Id, &newTrs.BillDate, &newTrs.EntryDate, &newTrs.FinishDate, &newTrs.EmployeeId, &newTrs.CustomerId)
	if errx != nil {
		tx.Rollback()
		fmt.Println(errx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Transaction"})
		return
	}

	// newTrs.BillDetails[0].Id = "1"
	for i, val := range newTrs.BillDetails {
		if !checkProductIdExist(val.ProductId) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product Id Tidak Ada"})
			return
		}
		if val.Qty <= 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Qty Kosong / Invalid"})
			return
		}
		errx := tx.QueryRow(queryInsertDetail, &newTrs.Id, &val.ProductId, &val.Qty).Scan(&val.Id, &val.BillId, &val.ProductId)
		if errx != nil {
			tx.Rollback()
			fmt.Println(errx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Transaction Detail"})
			return
		}
		price := tx.QueryRow(selectPrice, val.Id, val.ProductId).Scan(&val.ProductPrice)
		if price != nil {
			tx.Rollback()
			// fmt.Println(errx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Price Detail"})
			return
		}
		newTrs.BillDetails[i].Id = val.Id
		newTrs.BillDetails[i].BillId = val.BillId
		newTrs.BillDetails[i].ProductId = val.ProductId
		newTrs.BillDetails[i].ProductPrice = val.ProductPrice
		newTrs.BillDetails[i].Qty = val.Qty
	}

	errx = tx.Commit()
	if errx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Commit Input Transaksi " + errx.Error()})
		return
	}
	newResult.Message = "Input Berhasil"
	newResult.Data = newTrs
	c.JSON(http.StatusCreated, newResult)
}

func GetTransaksiByIdBill(c *gin.Context) {
	var idbil = c.Param("id_bill")
	query := `
			SELECT 
			aa.id as idtrs,
			TO_CHAR(aa.billdate,'YYYY-MM-DD HH24:MI:SS') as billdate,
			TO_CHAR(aa.entrydate,'YYYY-MM-DD HH24:MI:SS') as entrydate,
			TO_CHAR(aa.finishdate,'YYYY-MM-DD HH24:MI:SS') as finishdate,
			bb.id as empid,bb.name as nameemp,bb.phonenumber as phoneemp,bb.address as addremp,
			cc.id as csid, cc.name as namecus, cc.phonenumber as phonecs, cc.address as addrcs
			FROM trs_laundry as aa
			INNER JOIN 
			mst_employee as bb
			ON aa.employeeid = bb.id
			INNER JOIN 
			mst_customer as cc
			ON aa.customerid = cc.id
			WHERE aa.id = $1 `

	querydetail := `
				SELECT 
				aa.id,aa.billid,
				bb.id as pid,bb.name as pname,
				bb.price,bb.unit,
				bb.price as pdprice ,aa.qty
				FROM
				trs_laundry_detail aa
				INNER JOIN
				mst_product as bb
				ON aa.productid = bb.id
				WHERE aa.billid = $1`

	var dt TrsDetail
	tx, errrx := db.Begin()
	if errrx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + errrx.Error()})
		return
	}

	rows, err := tx.Query(query, idbil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error 1"})
		return
	}

	for rows.Next() {
		err := rows.Scan(&dt.Id, &dt.BillDate, &dt.EntryDate, &dt.FinishDate, &dt.Employee.Id, &dt.Employee.Name, &dt.Employee.PhoneNumber, &dt.Employee.Address, &dt.Customer.Id, &dt.Customer.Name, &dt.Customer.PhoneNumber, &dt.Customer.Address)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server Scan Detail error"})
			return
		}
	}

	rowsdetail, errr := tx.Query(querydetail, idbil)
	if errr != nil {
		fmt.Println(errr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error 2"})
		return
	}

	var bdt TrsBillDetail
	i := 0
	for rowsdetail.Next() {
		err := rowsdetail.Scan(&bdt.Id, &bdt.BillId, &bdt.Product.Id, &bdt.Product.Name, &bdt.Product.Price, &bdt.Product.Unit, &bdt.ProductPrice, &bdt.Qty)
		dt.BillDetails = append(dt.BillDetails, bdt)
		i++
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server Sacn Detail error 2"})
			return
		}
	}
	dt.TotalBiils = i

	errrx = tx.Commit()
	if errrx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Commit Input Transaksi " + errrx.Error()})
		return
	}
	var Result ResultDetail
	Result.Message = "Success"
	Result.Data = dt

	fmt.Println(len(dt.BillDetails))

	if dt.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id Not Found"})
	} else {
		c.JSON(http.StatusOK, Result)
	}
}

func GetTransaksi(c *gin.Context) {
	startDate := c.Query("startDate")
	layoutFormat := "02-01-2006 15:04:05"
	endDate := c.Query("endDate")
	productName := c.Query("productName")
	query := `
	SELECT 
	aa.id as idtrs,
	TO_CHAR(aa.billdate,'DD-MM-YYYY HH24:MI:SS') as billdate,
	TO_CHAR(aa.entrydate,'DD-MM-YYYY HH24:MI:SS') as entrydate,
	TO_CHAR(aa.finishdate,'DD-MM-YYYY HH24:MI:SS') as finishdate,
	bb.id as empid,bb.name as nameemp,bb.phonenumber as phoneemp,bb.address as addremp,
	cc.id as csid, cc.name as namecus, cc.phonenumber as phonecs, cc.address as addrcs
	FROM trs_laundry as aa
	INNER JOIN 
	mst_employee as bb
	ON aa.employeeid = bb.id
	INNER JOIN 
	mst_customer as cc
	ON aa.customerid = cc.id
	INNER JOIN
	trs_laundry_detail dd
	ON aa.id = dd.billid
	INNER JOIN
	mst_product as ee
	ON dd.productid = ee.id 
	`

	var rows *sql.Rows
	var err, errDate error
	var xstartDate, xendDate time.Time

	if startDate != "" && endDate != "" && productName != "" {
		xstartDate, errDate = time.Parse(layoutFormat, startDate)
		if errDate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Start Date DD-MM-YYYY H24:MI:SS"})
			return
		}
		xendDate, errDate = time.Parse(layoutFormat, endDate)
		if errDate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format End Date DD-MM-YYYY H24:MI:SS"})
			return
		}
		query += `
			 WHERE 
				aa.billdate BETWEEN $1 AND $2
				OR
				aa.finishdate BETWEEN $1 AND $2
				OR
				aa.entrydate BETWEEN $1 AND $2
				OR
				ee.name ILIKE '%'|| $3 ||'%' 
				ORDER BY aa.id DESC
		`
		rows, err = db.Query(query, xstartDate, xendDate, productName)
		fmt.Println("Kodisi 1")
	} else if startDate != "" && endDate != "" && productName == "" {
		xstartDate, errDate = time.Parse(layoutFormat, startDate)
		if errDate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Start Date DD-MM-YYYY H24:MI:SS"})
			return
		}
		xendDate, errDate = time.Parse(layoutFormat, endDate)
		if errDate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format End Date DD-MM-YYYY H24:MI:SS"})
			return
		}
		query += `
			 WHERE 
				aa.billdate BETWEEN $1 AND $2
				OR
				aa.finishdate BETWEEN $1 AND $2
				OR
				aa.entrydate BETWEEN $1 AND $2 
				ORDER BY aa.id DESC
		`
		rows, err = db.Query(query, xstartDate, xendDate)
		fmt.Println("Kodisi 2")

	} else if startDate == "" && endDate == "" && productName != "" {
		query += `
			 WHERE ee.name ILIKE '%'|| $1 ||'%' 
			 ORDER BY aa.id DESC
		`
		rows, err = db.Query(query, productName)
		fmt.Println("Kodisi 3")
	} else {
		query += " ORDER BY aa.id DESC "
		rows, err = db.Query(query)
		fmt.Println("Kodisi 4")
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Header server error"})
		return
	}

	var dt TrsDetail
	var Result ResultSearchDetail
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dt.Id, &dt.BillDate, &dt.EntryDate, &dt.FinishDate, &dt.Employee.Id, &dt.Employee.Name, &dt.Employee.PhoneNumber, &dt.Employee.Address, &dt.Customer.Id, &dt.Customer.Name, &dt.Customer.PhoneNumber, &dt.Customer.Address)
		dt.BillDetails = detailTransactionList(dt.Id)
		dt.TotalBiils = len(detailTransactionList(dt.Id))
		Result.Data = append(Result.Data, dt)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server Scan Detail error"})
			return
		}
	}

	Result.Message = "Success"

	if dt.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": " Data Kosong / Pencarian Belum Sesuai"})
	} else {
		c.JSON(http.StatusOK, Result)
	}
}

func detailTransactionList(idbill string) []TrsBillDetail {
	queryDetail := `
				SELECT
				bb.id as idtrsdetail,bb.billid,
				cc.id as pid,cc.name as pname,
				cc.price,cc.unit,
				cc.price as pdprice,bb.qty
				FROM trs_laundry as aa
				INNER JOIN
				trs_laundry_detail bb
				ON aa.id = bb.billid
				INNER JOIN
				mst_product as cc
				ON bb.productid = cc.id 
				WHERE aa.id = $1
				`
	rows, err := db.Query(queryDetail, idbill)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var xdt []TrsBillDetail
	var bdt TrsBillDetail
	for rows.Next() {
		err = rows.Scan(&bdt.Id, &bdt.BillId, &bdt.Product.Id, &bdt.Product.Name, &bdt.Product.Price, &bdt.Product.Unit, &bdt.ProductPrice, &bdt.Qty)
		if err != nil {
			fmt.Println(err)
		}
		xdt = append(xdt, bdt)
	}
	return xdt
}

func createTrsId() string {
	var result string
	date := getdate.GetYMD()
	var trsId string
	selectTrsId := "SELECT id FROM trs_laundry WHERE TO_CHAR(billdate,'YYYY-MM-DD') ILIKE '%' || $1 || '%' ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(selectTrsId, getdate.Get_YMD()).Scan(&trsId)
	if err != nil {
		result = "TRS" + date + "000001"
	} else {
		xtrsId := trsId[11:]
		csInt, _ := strconv.Atoi(xtrsId)
		csInt++
		lenindex := len(xtrsId) - len(strconv.Itoa(csInt))
		zero := ""
		for i := 0; i < lenindex; i++ {
			zero += "0"
		}
		// result = trsId[0 : len(trsId)-len(xtrsId)]
		result += "TRS" + date + zero + strconv.Itoa(csInt)
	}
	return result
}

func checkDate(dateInput string) bool {
	var layoutTime = "2006-01-02 15:04:05"
	_, err := time.Parse(layoutTime, dateInput)
	if err != nil {
		return false
	} else {
		return true
	}
}

func checkEmployeeIdExist(employeeId string) bool {
	var result bool = false
	selectEmployee := "SELECT id FROM mst_employee WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(selectEmployee, employeeId)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}

func checkCustomerIdExist(customerId string) bool {
	var result bool = false
	selectCustomer := "SELECT id FROM mst_customer WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(selectCustomer, customerId)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}

func checkProductIdExist(productId string) bool {
	var result bool = false
	selectProduct := "SELECT id FROM mst_product WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(selectProduct, productId)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}
