package helper

import (
	"data"
	"database/sql"
	"fmt"
	"log"
)

func AddStockIfNotExist(id int32, name string, db *sql.DB) {
	result := db.QueryRow("select * from `tbl_base` where `stockid`=%d", id)
	stockid := 0
	stock_name := ""
	err := result.Scan(&stockid, &stock_name)
	if err != nil {
		stmt, ok := db.Prepare("insert into `tbl_base` (`stockid`,`stockname`) values (?,?)")
		if ok == nil {
			_, err := stmt.Exec(id, name)
			if err == nil {
				fmt.Println(id, "add success")
			}
		}
		stmt.Close()
	}
}

func CreateFinanceTable(date int, db *sql.DB) {
	/*create_sql := "create table `tbl_finance_%d` (`stockid` int(11) NOT NULL,
	  PRIMARY KEY (`stockid`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
		stmt, ok := db.Prepare(create_sql)
		if ok == nil {
			_, err := stmt.Exec(date)
			if err == nil {
				fmt.Println(id, "create table success")
			}
		}*/
}

func AddFinance(date int) {

}

func AddShareHolder(stockid int, date int32, num int32, db *sql.DB) {
	stmt, ok := db.Prepare("insert into `tbl_shareholder` (`stockid`,`date`,`num`) values (?,?,?)")
	if ok == nil {
		_, err := stmt.Exec(stockid, date, num)
		if err == nil {
			fmt.Println(stockid, date, "add success")
		} else {
			fmt.Println(stockid, date, "add failed")
		}
	} else {
		log.Fatal(ok)
	}
	stmt.Close()
}

func FindShareHolder(stockid int, results []data.ShareHolder, db *sql.DB) []data.ShareHolder {
	rows, ok := db.Query("select `date`,`num` from `tbl_shareholder` where `stockid`=? and `date`>=20181231 and `date`<=20190630 order by `date` desc", stockid)
	if ok == nil {
		for rows.Next() {
			var stock data.ShareHolder
			if err := rows.Scan(&stock.Date, &stock.Num); err != nil {
				fmt.Println("get stockid failed", stockid)
				continue
			}
			results = append(results, stock)
		}
	} else {
		log.Fatal(ok)
	}
	rows.Close()
	return results
}
