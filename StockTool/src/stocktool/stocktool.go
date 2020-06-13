package main

import (
	"database/sql"
	"encoding/json"
	"filters"
	"fmt"
	//"helper"
	"io/ioutil"
	"mysql"
	"net/http"
	"regexp"
	"runtime"
	"shareholder"
	"strconv"
	"strings"
)

var dbpath string = "root:19880110@tcp(127.0.0.1:8889)/stock?charset=utf8"
var db *sql.DB
var chl chan int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	db = mysql.Connect(dbpath)
	if db != nil {
		defer db.Close()
	}
	chl = make(chan int)
	// fmt.Printf("%s", "hh")
	for i := 0; i < 70; i++ {
		go doYearFinance(20200331, i+1)
		//go updateShareHolder()
		//go findShareHolderDecrease()
	}
	// var arr [2]string = [2]string{"22ss", "aaa"}
	// arrstr, _ := json.Marshal(arr)
	// fmt.Println(string(arrstr))
	for i := 0; i < 70; i++ {
		<-chl
		//fmt.Println(v, "end", i)
	}

	return
}

func findShareHolderDecrease() {
	defer doYearFinanceEnd(1)
	stockids := []int{}
	stockids = mysql.GetAllStock(db, stockids)
	l := len(stockids)
	for i := 0; i < l; i++ {
		istrue := filters.FindShareHolderDecrease(stockids[i], db)
		if istrue {
			fmt.Println(stockids[i], "shareholder decrease too much")
		}
	}
}

func updateShareHolder() {
	defer doYearFinanceEnd(1)
	stockids := []int{}
	stockids = mysql.GetAllStock(db, stockids)
	l := len(stockids)
	for i := 0; i < l; i++ {
		shareholder.UpdateShareHolder(stockids[i], db)
	}
}

func doYearFinanceEnd(page int) {
	//	fmt.Println(page, "endnow")
	chl <- page
}

//年报
func doYearFinance(date int, page int) {
	defer doYearFinanceEnd(page)
	url := ""
	url = fmt.Sprintf("http://datainterface.eastmoney.com/EM_DataCenter/JS.aspx?type=SR&sty=YJBB&fd=2020-03-31&st=13&sr=-1&p=%d&ps=60&js={pages:(pc),data:[(x)]}&stat=0&rt=49756324", page)
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	datastr := string(data)
	//fmt.Println(datastr)
	var reg *regexp.Regexp
	reg, err = regexp.Compile("\\[.*\\]")
	if err != nil {
		fmt.Println(err)
		return
	}
	datastr = reg.FindString(datastr)
	//fmt.Println(datastr)

	var datarr []string
	var temparr []string
	var stockid int64
	//var lirun int64
	fdata := json.Unmarshal([]byte(datastr), &datarr)
	if fdata == nil {
		for _, v := range datarr {
			temparr = strings.Split(v, ",")
			stockid, _ = strconv.ParseInt(temparr[0], 10, 32)
			//shouru, _ := strconv.ParseFloat(temparr[4], 32)
			//stockname := temparr[1]
			//helper.AddStockIfNotExist(int32(stockid), stockname, db)
			//lirun, _ := strconv.ParseFloat(temparr[8], 32)
			//income, _ := strconv.ParseFloat(temparr[5], 32)
			//perincome, _ := strconv.ParseFloat(temparr[2], 32)
			//maoli, _ := strconv.ParseFloat(temparr[13], 32)
			//percash, _ := strconv.ParseFloat(temparr[12], 32)
			lirun2, _ := strconv.ParseFloat(temparr[7], 32)
			//if lirun >= 20 && income > 10 && perincome > 0.01 && maoli > 30 && percash > 0 {
			if lirun2 >= 64000000 && lirun2 <= 65000000 {
				fmt.Println(stockid, temparr)
			}

			//lirun, _ = strconv.ParseInt(temparr[7], 10, 32)
		}

	}
}

func doDayInfo() {
	resp, err := http.Get("http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?type=CT&cmd=C.2&sty=FCOIATA&sortType=C&sortRule=-1&page=1&pageSize=20&js=var%20quote_123%3d{rank:[(x)],pages:(pc)}&token=7bc05d0d4c3c22ef9fca8c2a912d779c&jsName=quote_123&_g=0.7556139376967295")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	datastr := string(data)
	var reg *regexp.Regexp
	reg, err = regexp.Compile("\\[.*\\]")
	if err != nil {
		fmt.Println(err)
		return
	}
	datastr = reg.FindString(datastr)
	// fmt.Println(datastr)

	var datarr []string
	var temparr []string
	fdata := json.Unmarshal([]byte(datastr), &datarr)
	if fdata == nil {
		for _, v := range datarr {
			temparr = strings.Split(v, ",")
			fmt.Println(temparr)
		}

	}
}
