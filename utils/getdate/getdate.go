package getdate

import (
	"challenge-goapi/config"
)

var db = config.ConnectDB()

func GetYMD() string {
	var errx error
	sqlStatement := "SELECT TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD');"
	var getdatetime string
	errx = db.QueryRow(sqlStatement).Scan(&getdatetime)
	if errx != nil {
		panic(errx)
	}
	return getdatetime
}

func Get_YMD() string {
	var errx error
	sqlStatement := "SELECT TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD');"
	var getdatetime string
	errx = db.QueryRow(sqlStatement).Scan(&getdatetime)
	if errx != nil {
		panic(errx)
	}
	return getdatetime
}
