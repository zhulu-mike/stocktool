package shareholder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"helper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UpdateShareHolder(stockid int, db *sql.DB) {
	url := "http://f10.eastmoney.com/ShareholderResearch/ShareholderResearchAjax?code=%s"
	stock_code := ""
	if stockid < 400000 {
		stock_code = fmt.Sprintf("%s%06d", "SZ", stockid)
	} else {
		stock_code = fmt.Sprintf("%s%06d", "SH", stockid)
	}

	url = fmt.Sprintf(url, stock_code)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(string(data))
	data_str := string(data)
	data_reader := strings.NewReader(data_str)
	data_decoder := json.NewDecoder(data_reader)
	for i := 0; i < 3; i++ {
		t, err := data_decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v\n", t, t)
	}
	type Item struct {
		Rq, Gdrs, Gdrs_jsqbh, Rjltg, Rjltg_jsqbh, cmjzd, gj, rjcgje, qsdgdcghj, qsdltgdcghj string
	}
	for data_decoder.More() {
		var item Item
		err := data_decoder.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		rq_time, _ := time.Parse("2006-01-02", item.Rq)
		date64, _ := strconv.ParseInt(rq_time.Format("20060102"), 10, 32)
		date32 := int32(date64)
		is_wan := strings.Index(item.Gdrs, "万")
		num64, _ := strconv.ParseFloat(strings.TrimSuffix(item.Gdrs, "万"), 32)
		num32 := int32(num64 * 10000)
		if is_wan == -1 {
			num32 = int32(num64)
		}
		fmt.Println(item.Rq, item.Gdrs, date32, num32)
		helper.AddShareHolder(stockid, date32, num32, db)
	}

}
