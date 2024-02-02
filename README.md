# Backend Laundry GO API 
### Persiapan :
1. Buat terlebih dahulu Database di PostgreSQl dengan Nama "tesLaundryGoApi"
2. Jalankan Script DDL yang ada pada file "/query/DDL.sql"
3. Untuk konfigurasi koneksi database bisa di sesuaikan pada file "/config/config.go"
4. Running aplikasi dengan mengetikan perintah "go run ." pada terminal / console
=====================================================================================
### Detail Menu :
1. Menu Customer
* Terdiri dari : id,name,phoneNumber,address.

**a. Get Customer By Id**
<br/>
**Url : /customers/:id**
 - akan tampil pesan "Id Not Found" jika Id yang di masukan tidak ada
 - parameter id yang di cari tidak case sensitive, hasilnya akan tetap sama baik dengan huruf kapital atau huruf kecil
<br/>

**b. Post Data Customer**
<br/>
**Url : /customers**
 - id(string) : terbuat otomatis dari function yang sudah di siapkan dengan mengambil tahun,bulan dan tahun realtime dikombinasikan dengan kode customer dan counter
 - name(string) : isi nama tidak boleh kosong, jika kosong akan tampil warning "Nama Customer Kosong"
 - phoneNumber(string) : isi nomor telepon atau handphone tidak boleh kosong, jika kosong akan tampil warning "Telepon Customer Kosong" dan harus unik / belum pernah di masukan di database, jika nomer telepon sama dengan yang sudah ada di database maka akan mereturn response "Telepon Customer Sudah Ada !"
 - address(string) : isi alamat tidak boleng kosong dan jumlah karakter minimal 10 karakter yang di masukan, ketika syarat yang ditentukan tidak sesuai maka akan merturn response "Alamat Customer Kosong" atau "Alamat Customer Minimal 10 Karakter"
<br/>

 **c. Update Data Customer**
 <br/>
 **Url : /customers/:id**
 - Sebelum bisa melakukan update ada pengecekan bahwa id customer yang dimasukan ada didalam database.
 Data akan terupdate dan terisi sesuai dengan yang di input, jika data yang di input kosong maka component / row di database tidak terupdate, hanya akan mengupdate row yang input nilainya tidak kosong
<br/>

**d. Delete Data Customer**
<br/>
**Url : /customers/:id**
- Sebelum bisa melakukan delete ada pengecekan bahwa id customer yang dimasukan ada didalam database. Jikad id customer tidak ada maka akan mereturn response "Customer Id Tidak Ada" dan tidak ada aksi delete.
--------------------------------------------------------------------------------------

2. Product
* Terdiri dari : id,name,price,unit.

**a. Get All Product**
<br/>

**Url : /products**
- Menampilkan semua data yang ada pada table mst_product
<br/>

**b. Get Product By Id**
<br/>
**Url : /products/:id**
 - akan tampil pesan "Id Not Found" jika Id yang di masukan tidak ada
 - parameter id yang di cari tidak case sensitive, hasilnya akan tetap sama baik dengan huruf kapital atau huruf kecil
<br/>

**c. Create Product**
<br/>
**Url : /products**
- id(string) : terbuat otomatis dari function yang sudah di siapkan dengan mengambil tahun,bulan dan tahun realtime dikombinasikan dengan kode product dan counter
- name(string) : isi nama tidak boleh kosong, jika kosong akan tampil warning "Nama Product Kosong"
- price(integer) : input price tidak boleh kosong dan harus angka yang bernilai lebih besar dari nol, jika tidak maka akan tampil warning "Harga Product Invalid"
- unit : input unit atau satuan barang atau service tidak boleh kosong, jika kosong akan tampil warning "Satuan Product Koson"
<br/>

**d. Update Product By Id**
<br/>
**Url : /products/:id**
- Sebelum bisa melakukan update ada pengecekan bahwa id product yang dimasukan ada didalam database.
Data akan terupdate dan terisi sesuai dengan yang di input, jika data yang di input kosong maka component / row di database tidak terupdate, hanya akan mengupdate row yang input nilainya tidak kosong
<br/>

**e. Delete Product By Id**
<br/>
**Url : /products/:id**
- Sebelum bisa melakukan delete ada pengecekan bahwa id product yang dimasukan ada didalam database. Jikad id product tidak ada maka akan mereturn response "Product Id Tidak Ada" dan tidak ada aksi delete.
--------------------------------------------------------------------------------------

3. Menu Employee
* Terdiri dari : id,name,phoneNumber,address.

**a. Get Employee By Id**
<br/>
**Url : /employees/:id**
 - akan tampil pesan "Id Not Found" jika Id yang di masukan tidak ada
 - parameter id yang di cari tidak case sensitive, hasilnya akan tetap sama baik dengan huruf kapital atau huruf kecil
<br/>

**b. Post Data Employee**
<br/>
**Url : /employees**
 - id(string) : terbuat otomatis dari function yang sudah di siapkan dengan mengambil tahun,bulan dan tahun realtime dikombinasikan dengan kode employee dan counter
 - name(string) : isi nama tidak boleh kosong, jika kosong akan tampil warning "Nama Employee Kosong"
 - phoneNumber(string) : isi nomor telepon atau handphone tidak boleh kosong, jika kosong akan tampil warning "Telepon Employee Kosong" dan harus unik / belum pernah di masukan di database, jika nomer telepon sama dengan yang sudah ada di database maka akan mereturn response "Telepon Employee Sudah Ada !"
 - address(string) : isi alamat tidak boleng kosong dan jumlah karakter minimal 10 karakter yang di masukan, ketika syarat yang ditentukan tidak sesuai maka akan merturn response "Alamat Employee Kosong" atau "Alamat Employee Minimal 10 Karakter"
<br/>

 **c. Update Data Employee**
 <br/>
 **Url : /employees/:id**
 - Sebelum bisa melakukan update ada pengecekan bahwa id employee yang dimasukan ada didalam database.
 Data akan terupdate dan terisi sesuai dengan yang di input, jika data yang di input kosong maka component / row di database tidak terupdate, hanya akan mengupdate row yang input nilainya tidak kosong
<br/>

**d. Delete Data Employee**
<br/>
**Url : /employees/:id**
- Sebelum bisa melakukan delete ada pengecekan bahwa id employee yang dimasukan ada didalam database. Jika id employee tidak ada maka akan mereturn response "Employee Id Tidak Ada" dan tidak ada aksi delete.
-------------------------------------------------------------------------------------

4. Transaksi
* Terdiri dari : "billDate,entryDate,finishDate,employeeId,customerId,billDetails.productId,billDetails.qty"

**a. POST Transaksi Laundry**
<br/>
**Url : /transactions**
- billDate(string) : format tanggal mengikuti deafult format tanggal PostgreSQL YYYY-MM-DD HH24:MI:SS. Contoh  : 2024-01-01 19:00:00
- entryDate(string) : format tanggal mengikuti deafult format tanggal PostgreSQL YYYY-MM-DD HH24:MI:SS. Contoh  : 2024-01-01 19:00:00
- finishDate(string) : format tanggal mengikuti deafult format tanggal PostgreSQL YYYY-MM-DD HH24:MI:SS. Contoh  : 2024-01-01 19:00:00
- employeeId(string) : id employee yang ada di tabel mst_employee 
- customerId(string) : id customer yang ada di tabel mst_customer
- billDetails.productId(string) : id product yang ada di tabel mst_product
- billDetails.qty(integer) : jumlah ataua kuantitas jasa / barang yang akan di transaksikan, input qty tidak boleh lebih kecil atau sama dengan nol
<br/>

**b. GET Transaksi Laundry By Id**
<br/>
**Url : /transactions/:id_bill**
- id(string) : merupakan id transaksi yang sudah ada di tabel trs_laundry
-  Menampilkan detail transaksi dari table trs_laundry, trs_laundry_detail, mst_product, mst_customer,mst_employee berdasarkan id transaksi yang di input
<br/>

**c. GET Transaksi Laundry By Params**
<br/>
**Url : /transactions**
- startDate(string) : format tanggal yaitu DD-MM-YY H24-M1-SS jika format tanggal tidak sesuai makan akan mengirip response error
- endDate(string) : format tanggal yaitu DD-MM-YY H24-M1-SS jika format tanggal tidak sesuai makan akan mengirip response error
- productName(string) : keyword atau kata kunci pencarian yang bisa di masukan berdasarkan nama product yang sebelumnya sudah di transaksikan
- ketiga parameter diatas optional, tapi untuk start date dan end date harus terisi secara bersamaan, tidak dianjurkan jika hanya menggunakan startdate atau enddate saja
contoh : /transactions?starDate=01-01-2024 00:00:00&endDate=31-01-2024 23:59:59&productName=CUCI
<br/>
<br/>