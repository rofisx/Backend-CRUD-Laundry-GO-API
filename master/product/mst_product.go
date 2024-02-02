package product

import (
	"challenge-goapi/config"
	"challenge-goapi/utils/getdate"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Unit  string `json:"unit"`
}
type ResultAllProduct struct {
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}
type ResultProductById struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}
type Result struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

var db = config.ConnectDB()

func GetAllProduct(c *gin.Context) {
	query := "SELECT id,name,price,unit FROM mst_product"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		products = append(products, p)
	}
	var res ResultAllProduct
	res.Message = "Success"
	res.Data = products
	c.JSON(http.StatusOK, res)
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id,name,price,unit FROM mst_product WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(query, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()
	var p Product
	for rows.Next() {
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	if p.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id Product Not Found"})
	} else {
		var res ResultProductById
		res.Message = "Success"
		res.Data = p
		c.JSON(http.StatusOK, res)
	}
}

func CreateProduct(c *gin.Context) {
	var newP Product
	err := c.ShouldBind(&newP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newP.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Product Kosong"})
		return
	} else {
		if checkNamaProductExist(newP.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Product Sudah Ada"})
			return
		}
	}
	if newP.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga Product Invalid"})
		return
	}

	if newP.Unit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Satuan Product Kosong"})
		return
	}

	queryInsert := "INSERT INTO mst_product (id,name,price,unit) VALUES ($1,$2,$3,$4) RETURNING id"
	var prdId string
	newP.Id = createProductId()
	fmt.Println(newP.Id)
	err = db.QueryRow(queryInsert, newP.Id, newP.Name, newP.Price, newP.Unit).Scan(&prdId)
	if err != nil {
		// fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to created Product"})
		return
	}
	newP.Id = prdId
	var res ResultProductById
	res.Message = "Input Berhasil"
	res.Data = newP
	c.JSON(http.StatusCreated, res)
}

func UpdateProductById(c *gin.Context) {
	id := c.Param("id")
	var updPrd Product
	var valPrd Product
	err := c.ShouldBind(&updPrd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := "SELECT id,name,price,unit FROM mst_product WHERE LOWER(id) = LOWER($1)"
	rows, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	for rows.Next() {
		err = rows.Scan(&valPrd.Id, &valPrd.Name, &valPrd.Price, &valPrd.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	// c.JSON(http.StatusOK, valCs)
	defer rows.Close()

	if len(valPrd.Id) > 0 {
		updPrd.Id = valPrd.Id
		if updPrd.Name == "" {
			updPrd.Name = valPrd.Name
		}
		if updPrd.Price <= 0 {
			updPrd.Price = valPrd.Price
		}

		if updPrd.Unit == "" {
			updPrd.Unit = valPrd.Unit
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		updQuery := "UPDATE mst_product SET name = $2, price = $3, unit = $4 WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(updQuery, updPrd.Id, updPrd.Name, updPrd.Price, updPrd.Unit)
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
			var res ResultProductById
			res.Message = "Update Berhasil"
			res.Data = updPrd
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product Id Tidak Ada"})
		return
	}
}

func DeleteProductById(c *gin.Context) {
	id := c.Param("id")
	var idprd string
	query := "SELECT id FROM mst_product WHERE LOWER(id) = LOWER($1)"
	err := db.QueryRow(query, id).Scan(&idprd)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product Id Tidak Ada"})
		return
	}

	if len(idprd) > 0 {
		querytrscheck := "SELECT productid FROM trs_laundry_detail WHERE LOWER(productid) = LOWER($1)"
		errtrs := db.QueryRow(querytrscheck, id).Scan(&idprd)
		// fmt.Println("error check", errtrs)
		if errtrs == nil {
			fmt.Println(errtrs)
			c.JSON(http.StatusBadRequest, gin.H{"error": "product id tidak bisa dihapus, sudah digunakan di transaksi"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Gaes" + err.Error()})
			return
		}

		delQuery := "DELETE FROM mst_product WHERE LOWER(id) = LOWER($1)"
		_, errx := tx.Exec(delQuery, idprd)
		if errx != nil {
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product Id Tidak Ada"})
		return
	}
}

func createProductId() string {
	var result string
	date := getdate.GetYMD()
	var prdId string
	selectProduct := "SELECT id FROM mst_product ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(selectProduct).Scan(&prdId)
	if err != nil {
		result = "SERV" + date + "000001"
	} else {
		xprdId := prdId[12:]
		prdIdInt, _ := strconv.Atoi(xprdId)
		prdIdInt++
		lenindex := len(xprdId) - len(strconv.Itoa(prdIdInt))
		zero := ""
		for i := 0; i < lenindex; i++ {
			zero += "0"
		}
		// result = prdId[0 : len(prdId)-len(xprdId)]
		result += "SERV" + date + zero + strconv.Itoa(prdIdInt)
	}
	return result
}

func checkNamaProductExist(product string) bool {
	var result bool = false
	selectProduct := "SELECT name FROM mst_product WHERE LOWER(name) = LOWER($1)"
	rows, err := db.Query(selectProduct, product)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		result = true
	}
	return result
}
