-- CREATE Master Customer
CREATE TABLE mst_customer(
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	phoneNumber VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL
);

-- CREATE Master product
CREATE TABLE mst_product(
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	price INT NOT NULL,
	unit VARCHAR(100) NOT NULL
)

-- CREATE Master employee
CREATE TABLE mst_employee(
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	phoneNumber VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL
);

-- CREATE Transaksi Header Laundry
CREATE TABLE trs_laundry(
	id VARCHAR(255) PRIMARY KEY,
	billDate TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	entryDate TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	finishDate TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	employeeId VARCHAR(255) NOT NULL,
	customerId VARCHAR(255) NOT NULL,
	FOREIGN KEY (employeeId) REFERENCES mst_employee(id),
	FOREIGN KEY (customerId) REFERENCES mst_customer(id)
);

-- CREATE Transaksi Detail Laundry
CREATE TABLE trs_laundry_detail(
	id SERIAL PRIMARY KEY,
	billId VARCHAR(255) NOT NULL,
	productId VARCHAR(255) NOT NULL,
	qty int,
	FOREIGN KEY (billId) REFERENCES trs_laundry(id),
	FOREIGN KEY (productId) REFERENCES mst_product(id)
);