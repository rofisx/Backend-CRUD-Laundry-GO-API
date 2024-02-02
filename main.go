package main

import (
	"challenge-goapi/master/customer"
	"challenge-goapi/master/employee"
	"challenge-goapi/master/product"
	"challenge-goapi/master/transaction"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Get All Customer
	router.GET("/customers", customer.GetAllCustomer)
	// Get Customer By Id
	router.GET("/customers/:id", customer.GetCustomerById)
	//POST Cuctomer
	router.POST("/customers", customer.CreateCustomer)
	//PUT Cuctomer
	router.PUT("/customers/:id", customer.UpdateCustomerById)
	//DELETE Cuctomer
	router.DELETE("/customers/:id", customer.DeleteCustomerById)

	//Get All Products
	router.GET("/products", product.GetAllProduct)
	//Get Product By Id
	router.GET("/products/:id", product.GetProductById)
	//POST Product
	router.POST("/products", product.CreateProduct)
	//PUT Product By Id
	router.PUT("/products/:id", product.UpdateProductById)
	//DELETE Product By Id
	router.DELETE("/products/:id", product.DeleteProductById)

	// Get All Customer
	router.GET("/employees", employee.GetAllEmployee)
	// Get Customer By Id
	router.GET("/employees/:id", employee.GetEmployeeById)
	//POST Cuctomer
	router.POST("/employees", employee.CreateEmployee)
	//PUT Cuctomer
	router.PUT("/employees/:id", employee.UpdateEmployeeById)
	//DELETE Cuctomer
	router.DELETE("/employees/:id", employee.DeleteEmployeeById)

	//POST Transaksi
	router.POST("/transactions", transaction.CreateTransaksiLaundry)

	//GET Transaksi
	router.GET("/transactions/:id_bill", transaction.GetTransaksiByIdBill)
	//GET All Transaksi
	router.GET("/transactions", transaction.GetTransaksi)

	router.Run(":8080")
}
