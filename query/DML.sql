/*=== Customer ===*/
-- INSERT mst_customer
"INSERT INTO mst_customer (id,name,phonenumber,address) VALUES ('CUST20240125000001','Customer 1','08674672346324','Jl. Cebolok Semarang');"

-- UPDATE mst_customer
"UPDATE mst_customer SET name = 'Customer 2', phonenumber = '08232365987065',address = 'Jl.Cebolok No.6 Semarang' WHERE id = 'CUST20240125000001';"

-- DELETE mst_customer
"DELETE FROM mst_customer WHERE id = 'CUST20240125000001';"

/*=== Product ===*/
-- INSERT mst_product
"INSERT INTO mst_product (id,name,price,unit) VALUES ('SERV20240125000001','CUCI KERING',8000,'PCS');"

-- UPDATE mst_product
"UPDATE mst_product SET name = 'CUCI + SETRIKA', price = 10000, unit = 'PCS' WHERE id = 'SERV20240125000001' "

-- DELETE mst_product
"DELETE FROM mst_product WHERE id = 'SERV20240125000001' "

/*=== Employee ===*/
-- INSERT mst_employee
"INSERT INTO mst_employee (id,name,phonenumber,address) VALUES ('EMP20240126000001','Employee 1','08674672346324','Jl. Cebolok Semarang');"

-- UPDATE mst_employee
"UPDATE mst_employee SET name = 'employee 2', phonenumber = '08232365987065',address = 'Jl.Cebolok No.6 Semarang' WHERE id = 'EMP20240126000001';"

-- DELETE mst_employee
"DELETE FROM mst_employee WHERE id = 'EMP20240126000001';"


/*=== Transaction ===*/
INSERT INTO trs_laundry(
id, billdate, entrydate, finishdate, employeeid, customerid)
VALUES ('TRS20240202000001', '2024-01-31 07:00:00', '2024-01-31 07:00:00', '2024-01-31 10:00:00', 'EMP20240126000001', 'CUST20240124000002');

/*=== Transaction Detail ===*/
INSERT INTO trs_laundry_detail(
billid, productid, qty)
VALUES ( 'TRS20240202000001', 'SERV20240125000001', 2);