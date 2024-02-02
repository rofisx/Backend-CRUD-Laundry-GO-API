package employee

import (
	"challenge-goapi/config"
	"challenge-goapi/utils/getdate"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type ResultAllEmployee struct {
	Message string     `json:"message"`
	Data    []Employee `json:"data"`
}
type ResultEmployeeById struct {
	Message string   `json:"message"`
	Data    Employee `json:"data"`
}
type Result struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

var db = config.ConnectDB()

func GetAllEmployee(c *gin.Context) {
	query := "SELECT id,name,phonenumber,address FROM mst_employee"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()
	var empe []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		empe = append(empe, emp)
	}
	if len(empe) > 0 {
		var res ResultAllEmployee
		res.Message = "Success"
		res.Data = empe
		c.JSON(http.StatusOK, res)
	} else {
		var xres Result
		xres.Message = "Success"
		xres.Data = "Belum ada data employee"
		c.JSON(http.StatusOK, xres)
	}
}

func GetEmployeeById(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id,name,phonenumber,address FROM mst_employee WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()
	var emp Employee
	for rows.Next() {
		err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	var res ResultEmployeeById
	if emp.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id Not Found"})
	} else {
		res.Message = "Success"
		res.Data = emp
		c.JSON(http.StatusOK, res)
	}
}

func CreateEmployee(c *gin.Context) {
	var newEmp Employee
	err := c.ShouldBind(&newEmp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newEmp.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Employee Kosong"})
		return
	}
	if newEmp.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telepon Employee Kosong"})
		return
	} else {
		if checkphoneEmployeeExist(newEmp.PhoneNumber) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Telepon Employee Sudah Ada !"})
			return
		}
	}
	if newEmp.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat Employee Kosong"})
		return
	} else {
		if len(newEmp.Address) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat Employee Minimal 10 Karakter"})
			return
		}
	}

	queryInsert := "INSERT INTO mst_employee (id,name,phonenumber,address) VALUES ($1,$2,$3,$4) RETURNING id"
	var empId string
	newEmp.Id = createEmployeeId()
	// fmt.Println("== ada ==", newEmp.Id)
	err = db.QueryRow(queryInsert, newEmp.Id, newEmp.Name, newEmp.PhoneNumber, newEmp.Address).Scan(&empId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Employee"})
		return
	}
	newEmp.Id = empId
	var res ResultEmployeeById
	res.Message = "Success"
	res.Data = newEmp
	c.JSON(http.StatusCreated, res)
}

func UpdateEmployeeById(c *gin.Context) {
	id := c.Param("id")
	var updEmp Employee
	var valEmp Employee
	err := c.ShouldBind(&updEmp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := "SELECT id,name,phonenumber,address FROM mst_employee WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	for rows.Next() {
		err = rows.Scan(&valEmp.Id, &valEmp.Name, &valEmp.PhoneNumber, &valEmp.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	// c.JSON(http.StatusOK, valEmp)
	defer rows.Close()

	if len(valEmp.Id) > 0 {
		updEmp.Id = valEmp.Id
		if updEmp.Name == "" {
			updEmp.Name = valEmp.Name
		}
		if updEmp.PhoneNumber == "" {
			updEmp.PhoneNumber = valEmp.PhoneNumber
		}

		if updEmp.Address == "" {
			updEmp.Address = valEmp.Address
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		updQuery := "UPDATE mst_employee SET name = $2, phonenumber = $3, address = $4 WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(updQuery, updEmp.Id, updEmp.Name, updEmp.PhoneNumber, updEmp.Address)
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
			var res ResultEmployeeById
			res.Message = "Success"
			res.Data = updEmp
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee Id Tidak Ada"})
		return
	}
}

func DeleteEmployeeById(c *gin.Context) {
	id := c.Param("id")
	var idemp string
	query := "SELECT id FROM mst_employee WHERE LOWER(id) = LOWER($1)"
	err := db.QueryRow(query, id).Scan(&idemp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee Id Tidak Ada"})
		return
	}

	if len(idemp) > 0 {
		querytrscheck := "SELECT employeeid FROM trs_laundry WHERE LOWER(employeeid) = LOWER($1)"
		errtrs := db.QueryRow(querytrscheck, id).Scan(&idemp)
		fmt.Println(errtrs)
		if errtrs == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "employee id tidak bisa dihapus, sudah digunakan di transaksi"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		delQuery := "DELETE FROM mst_employee WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(delQuery, idemp)
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

func createEmployeeId() string {
	var result string
	date := getdate.GetYMD()
	var csId string
	selectService := "SELECT id FROM mst_employee ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(selectService).Scan(&csId)
	if err != nil {
		result = "EMP" + date + "000001"
	} else {
		xcsId := csId[11:]
		csInt, _ := strconv.Atoi(xcsId)
		csInt++
		lenindex := len(xcsId) - len(strconv.Itoa(csInt))
		zero := ""
		for i := 0; i < lenindex; i++ {
			zero += "0"
		}
		// result = csId[0 : len(csId)-len(xcsId)]
		result += "EMP" + date + zero + strconv.Itoa(csInt)
	}
	return result
}

func checkphoneEmployeeExist(csphone string) bool {
	var result bool = false
	selectService := "SELECT phonenumber FROM mst_employee WHERE phonenumber = $1"
	rows, err := db.Query(selectService, csphone)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}
