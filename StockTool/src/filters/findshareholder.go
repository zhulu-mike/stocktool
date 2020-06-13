package filters

import (
	"data"
	"database/sql"
	"fmt"
	"helper"
)

func FindShareHolderDecrease(stockid int, db *sql.DB) bool {
	results := []data.ShareHolder{}
	results = helper.FindShareHolder(stockid, results, db)
	//
	l := len(results)
	if l == 0 {
		return false
	}
	min_num := results[l-1].Num
	max_num := results[0].Num
	scale := float32(max_num-min_num) / float32(min_num) * 100

	if int32(scale) <= -40 {
		fmt.Println(results)
		fmt.Println(min_num, max_num, scale)
		return true
	}

	return false
}
