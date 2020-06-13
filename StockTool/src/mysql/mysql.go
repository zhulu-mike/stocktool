package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(dbpath string) *sql.DB {
	db, err := sql.Open("mysql", dbpath)
	if err != nil {
		fmt.Println("connect failed")
		return nil
	}
	fmt.Println("connect success")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		fmt.Println("connect valid")
	}
	return db
}

func GetOne(db *sql.DB, sqlstr string) {

}

func GetAllStock(db *sql.DB, stockids []int) []int {
	rows, err := db.Query("select `stockid` from `tbl_base` where 1")
	if err != nil {
		fmt.Println("get stock failed", err.Error())
		return stockids
	}
	for rows.Next() {
		stockid := 0
		if err := rows.Scan(&stockid); err != nil {
			fmt.Println("get stockid failed")
			continue
		}
		stockids = append(stockids, stockid)
	}
	rows.Close()
	return stockids
}
