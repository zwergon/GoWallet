package utils

import "database/sql"

//DBInfo definition de database
type DBInfo struct {
	User string
	Pass string
	Name string
	Host string
}

//NewDB create ad DBInfo struct ("Always or anything else")
func NewDB(key string) *DBInfo {
	dbInfo := DBInfo{}
	if key == "Always" {
		dbInfo.User = "zwergon_pwallet"
		dbInfo.Pass = "5sz3yqhr"
		dbInfo.Name = "zwergon_pwallet"
		dbInfo.Host = "mysql-zwergon.alwaysdata.net"
	} else {
		dbInfo.User = "root"
		dbInfo.Pass = "root"
		dbInfo.Name = "pwallet"
		dbInfo.Host = "localhost"
	}

	return &dbInfo
}

//DbConn connect to database
func (dbInfo *DBInfo) DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	db, err := sql.Open(dbDriver, dbInfo.User+":"+dbInfo.Pass+"@tcp("+dbInfo.Host+":3306)/"+dbInfo.Name)
	if err != nil {
		panic(err.Error())
	}
	return db
}
