package customer

import (
	"challenge-goapi/config"
	"challenge-goapi/utils/getdate"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
type ResultAllCustomer struct {
	Message string     `json:"message"`
	Data    []Customer `json:"data"`
}
type ResultCustomerById struct {
	Message string   `json:"message"`
	Data    Customer `json:"data"`
}
type Result struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func GetAllCustomer(c *gin.Context) {
	query := "SELECT id,name,phonenumber,address FROM mst_customer"

	rows, err := config.ConnectDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()
	var cs []Customer
	for rows.Next() {
		var xcs Customer
		err := rows.Scan(&xcs.Id, &xcs.Name, &xcs.PhoneNumber, &xcs.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		cs = append(cs, xcs)
	}
	if len(cs) > 0 {
		var res ResultAllCustomer
		res.Message = "Success"
		res.Data = cs
		c.JSON(http.StatusOK, res)
	} else {
		var xres Result
		xres.Message = "Success"
		xres.Data = "Belum ada data customer"
		c.JSON(http.StatusOK, xres)
	}
}

func GetCustomerById(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id,name,phonenumber,address FROM mst_customer WHERE LOWER(id) = LOWER($1)"
	rows, err := config.ConnectDB().Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()
	var cs Customer
	for rows.Next() {
		err := rows.Scan(&cs.Id, &cs.Name, &cs.PhoneNumber, &cs.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	var res ResultCustomerById
	res.Message = "Success"
	res.Data = cs
	if cs.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id Not Found"})
	} else {
		c.JSON(http.StatusOK, cs)
	}
}

func CreateCustomer(c *gin.Context) {
	var newCs Customer
	err := c.ShouldBind(&newCs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newCs.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Customer Kosong"})
		return
	}
	if newCs.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telepon Customer Kosong"})
		return
	} else {
		if checkphoneCustomerExist(newCs.PhoneNumber) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Telepon Customer Sudah Ada !"})
			return
		}
	}
	if newCs.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat Customer Kosong"})
		return
	} else {
		if len(newCs.Address) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat Customer Minimal 10 Karakter"})
			return
		}
	}

	queryInsert := "INSERT INTO mst_customer (id,name,phonenumber,address) VALUES ($1,$2,$3,$4) RETURNING id"
	var csId string
	newCs.Id = createCustomerId()
	// fmt.Println("== ada ==", newCs.Id)
	err = config.ConnectDB().QueryRow(queryInsert, newCs.Id, newCs.Name, newCs.PhoneNumber, newCs.Address).Scan(&csId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Customer"})
		return
	}
	newCs.Id = csId
	var res ResultCustomerById
	res.Message = "Input Berhasil"
	res.Data = newCs
	c.JSON(http.StatusCreated, res)
}

func UpdateCustomerById(c *gin.Context) {
	id := c.Param("id")
	var updCs Customer
	var valCs Customer
	err := c.ShouldBind(&updCs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := "SELECT id,name,phonenumber,address FROM mst_customer WHERE LOWER(id) = LOWER($1)"
	rows, err := config.ConnectDB().Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	for rows.Next() {
		err = rows.Scan(&valCs.Id, &valCs.Name, &valCs.PhoneNumber, &valCs.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	// c.JSON(http.StatusOK, valCs)
	defer rows.Close()

	if len(valCs.Id) > 0 {
		updCs.Id = valCs.Id
		if updCs.Name == "" {
			updCs.Name = valCs.Name
		}
		if updCs.PhoneNumber == "" {
			updCs.PhoneNumber = valCs.PhoneNumber
		}

		if updCs.Address == "" {
			updCs.Address = valCs.Address
		}

		tx, err := config.ConnectDB().Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		updQuery := "UPDATE mst_customer SET name = $2, phonenumber = $3, address = $4 WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(updQuery, updCs.Id, updCs.Name, updCs.PhoneNumber, updCs.Address)
		if errx != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Exec Update" + err.Error()})
			return
		}
		errx = tx.Commit()
		if errx != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Commit Update" + err.Error()})
			return
		} else {
			var res ResultCustomerById
			res.Message = "Update Berhasil"
			res.Data = updCs
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer Id Tidak Ada"})
		return
	}
}

func DeleteCustomerById(c *gin.Context) {
	id := c.Param("id")
	// var valCs Customer
	var idcs string
	query := "SELECT id FROM mst_customer WHERE LOWER(id) = LOWER($1)"
	err := config.ConnectDB().QueryRow(query, id).Scan(&idcs)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer Id Tidak Ada"})
		return
	}

	if len(idcs) > 0 {
		querytrscheck := "SELECT customerid FROM trs_laundry WHERE LOWER(customerid) = LOWER($1)"
		errtrs := config.ConnectDB().QueryRow(querytrscheck, id).Scan(&idcs)
		// fmt.Println("error check", errtrs)
		if errtrs == nil {
			fmt.Println(errtrs)
			c.JSON(http.StatusBadRequest, gin.H{"error": "customer id tidak bisa dihapus, sudah digunakan di transaksi"})
			return
		}

		tx, err := config.ConnectDB().Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		delQuery := "DELETE FROM mst_customer WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(delQuery, idcs)
		if errx != nil {
			fmt.Println(errx)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Exec Delete" + err.Error()})
			return
		}
		errx = tx.Commit()
		if errx != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Commit Delete" + err.Error()})
			return
		} else {
			var res Result
			res.Message = "Delete Berhasil"
			res.Data = "OK"
			c.JSON(http.StatusOK, res)
		}
	}
}

func createCustomerId() string {
	var result string
	date := getdate.GetYMD()
	var csId string
	selectService := "SELECT id FROM mst_customer ORDER BY id DESC LIMIT 1"
	err := config.ConnectDB().QueryRow(selectService).Scan(&csId)
	if err != nil {
		result = "CUST" + date + "000001"
	} else {
		xcsId := csId[12:]
		csInt, _ := strconv.Atoi(xcsId)
		csInt++
		lenindex := len(xcsId) - len(strconv.Itoa(csInt))
		zero := ""
		for i := 0; i < lenindex; i++ {
			zero += "0"
		}
		// result = csId[0 : len(csId)-len(xcsId)]
		result += "CUST" + date + zero + strconv.Itoa(csInt)
	}
	return result
}

func checkphoneCustomerExist(csphone string) bool {
	var result bool = false
	selectService := "SELECT phonenumber FROM mst_customer WHERE phonenumber = $1"
	rows, err := config.ConnectDB().Query(selectService, csphone)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}
