package transaction

import (
	"challenge-goapi/master/customer"
	"challenge-goapi/master/employee"
)

type TrsHeader struct {
	Id          string       `json:"id"`
	BillDate    string       `json:"billDate"`
	EntryDate   string       `json:"entryDate"`
	FinishDate  string       `json:"finishDate"`
	EmployeeId  string       `json:"employeeId"`
	CustomerId  string       `json:"customerId"`
	BillDetails []BillDetail `json:"billdetails"`
}
type BillDetail struct {
	Id           string `json:"id"`
	BillId       string `json:"billId"`
	ProductId    string `json:"productId"`
	ProductPrice int    `json:"productPrice"`
	Qty          int    `json:"qty"`
}

type Result struct {
	Message string    `json:"message"`
	Data    TrsHeader `json:"data"`
}

type ResultDetail struct {
	Message string    `json:"message"`
	Data    TrsDetail `json:"data"`
}

type ResultSearchDetail struct {
	Message string      `json:"message"`
	Data    []TrsDetail `json:"data"`
}

type TrsDetail struct {
	Id          string            `json:"id"`
	BillDate    string            `json:"billDate"`
	EntryDate   string            `json:"entryDate"`
	FinishDate  string            `json:"finishDate"`
	Employee    employee.Employee `json:"employee"`
	Customer    customer.Customer `json:"customer"`
	BillDetails []TrsBillDetail   `json:"billdetails"`
	TotalBiils  int               `json:"totalBill"`
}

type TrsBillDetail struct {
	Id           string        `json:"id"`
	BillId       string        `json:"billId"`
	Product      ProductDetail `json:"product"`
	ProductPrice int           `json:"productPrice"`
	Qty          int           `json:"qty"`
}

type ProductDetail struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Unit  string `json:"unit"`
}
