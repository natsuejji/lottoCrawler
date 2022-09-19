package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	// "os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

var db *sql.DB

type LottoInfo struct {
	drawDate    string
	period      string
	lottoType   string
	lottoNumber string
}

func yearConvert(year string) (string, error) {
	ceYear, err := strconv.ParseInt(year, 10, 0)
	if err != nil {
		return "", err
	}
	ceYear = ceYear + 1911
	return strconv.Itoa(int(ceYear)), nil
}
func search3DNumber(twYear int, month int) []LottoInfo {
	three_url := "https://www.taiwanlottery.com.tw/Lotto/3D/history.aspx"
	var lottoDatas []LottoInfo
	var formData = map[string]string{
		"L3DControl_history1$chk":       "radYM",
		"L3DControl_history1$dropYear":  strconv.Itoa(int(twYear)),
		"L3DControl_history1$dropMonth": strconv.Itoa(int(month)),
		"L3DControl_history1$btnSubmit": "查詢",
	}

	c := colly.NewCollector()
	c.OnHTML("input[name=__VIEWSTATE]", func(h *colly.HTMLElement) {
		formData["__VIEWSTATE"] = h.Attr("value")
	})
	c.OnHTML("input[name=__VIEWSTATEGENERATOR]", func(h *colly.HTMLElement) {
		formData["__VIEWSTATEGENERATOR"] = h.Attr("value")
	})
	c.OnHTML("input[name=__EVENTVALIDATION]", func(h *colly.HTMLElement) {
		formData["__EVENTVALIDATION"] = h.Attr("value")
	})
	c.Visit(three_url)
	c.OnHTML("#right > table > tbody > tr > td > table > tbody > tr:nth-child(3) > td:nth-child(2) > p", func(h *colly.HTMLElement) {
		var date = strings.Split(h.Text[6:], "/")
		var newLottoInfo LottoInfo
		newLottoInfo.lottoType = "3D"
		twYear, err := yearConvert(date[0])
		if err != nil {
			fmt.Println(err)
		}
		// get period
		newLottoInfo.period = h.DOM.Parent().Prev().Text()
		// get draw date
		date[0] = twYear
		date[1] = fmt.Sprintf("%02s", date[1])
		date[2] = fmt.Sprintf("%02s", date[2])
		dateString := strings.Join(date, "-")
		const layout = "2006-01-02"
		stemp, err := time.Parse(layout, dateString)
		if err != nil {
			fmt.Println(err, stemp)
		}
		// get lotto number
		var content = ""
		h.DOM.Parent().Next().Each(func(i int, s *goquery.Selection) {
			content += s.Text()
		})
		newLottoInfo.drawDate = strconv.Itoa(int(stemp.Unix()))
		newLottoInfo.lottoNumber = content
		lottoDatas = append(lottoDatas, newLottoInfo)
	})

	err := c.Post(three_url, formData)

	if err != nil {
		fmt.Println(err)
	}

	return lottoDatas
}
func search4DNumber(twYear int, month int) []LottoInfo {
	url := "https://www.taiwanlottery.com.tw/Lotto/4D/history.aspx"

	var formData = map[string]string{
		"L4DControl_history1$chk":       "radYM",
		"L4DControl_history1$dropYear":  strconv.Itoa(int(twYear)),
		"L4DControl_history1$dropMonth": strconv.Itoa(int(month)),
		"L4DControl_history1$btnSubmit": "查詢",
	}
	var lottoDatas []LottoInfo
	c := colly.NewCollector()
	c.OnHTML("input[name=__VIEWSTATE]", func(h *colly.HTMLElement) {
		formData["__VIEWSTATE"] = h.Attr("value")
	})
	c.OnHTML("input[name=__EVENTTARGET]", func(h *colly.HTMLElement) {
		formData["__EVENTTARGET"] = h.Attr("value")
	})
	c.OnHTML("input[name=__EVENTARGUMENT]", func(h *colly.HTMLElement) {
		formData["__EVENTARGUMENT"] = h.Attr("value")
	})
	c.OnHTML("input[name=__LASTFOCUS]", func(h *colly.HTMLElement) {
		formData["__LASTFOCUS"] = h.Attr("value")
	})

	c.OnHTML("input[name=__VIEWSTATEGENERATOR]", func(h *colly.HTMLElement) {
		formData["__VIEWSTATEGENERATOR"] = h.Attr("value")
	})
	c.OnHTML("input[name=__EVENTVALIDATION]", func(h *colly.HTMLElement) {
		formData["__EVENTVALIDATION"] = h.Attr("value")
	})
	c.Visit(url)

	c.OnHTML("table > tbody > tr > td > div:nth-child(8) > table > tbody > tr:nth-child(3) > td:nth-child(2) > p", func(h *colly.HTMLElement) {
		var date = strings.Split(h.Text[6:], "/")
		var newLottoInfo LottoInfo
		newLottoInfo.lottoType = "4D"
		twYear, err := yearConvert(date[0])
		if err != nil {
			fmt.Println(err)
		}

		// get period
		newLottoInfo.period = h.DOM.Parent().Prev().Text()
		date[0] = twYear
		date[1] = fmt.Sprintf("%02s", date[1])
		date[2] = fmt.Sprintf("%02s", date[2])
		dateString := strings.Join(date, "-")
		const layout = "2006-01-02"
		stemp, err := time.Parse(layout, dateString)
		if err != nil {
			fmt.Println(err, stemp)
		}
		var content = ""
		h.DOM.Parent().Next().Each(func(i int, s *goquery.Selection) {
			content += s.Text()
		})
		stempStr := strconv.Itoa(int(stemp.Unix()))
		newLottoInfo.drawDate = stempStr
		newLottoInfo.lottoNumber = content
		lottoDatas = append(lottoDatas, newLottoInfo)
	})

	err := c.Post(url, formData)

	if err != nil {
		fmt.Println(err)
	}
	return lottoDatas
}
func lottoDiaryTask3D() {

	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	c := colly.NewCollector()
	var date string
	var period string
	var number1 string
	var number2 string
	var number3 string
	c.OnHTML("div.contents_box04 > div#contents_logo_08", func(h *colly.HTMLElement) {
		date = strings.Split(h.DOM.Next().First().Text(), "\u00a0")[0]
		rex := regexp.MustCompile("[0-9]+")
		period = rex.FindString(strings.Split(h.DOM.Next().First().Text(), "\u00a0")[1])
		number := h.DOM.NextAll().Filter(".ball_tx").Text()
		number1 = number[0:1]
		number2 = number[1:2]
		number3 = number[2:3]
	})
	c.Visit(url)
	statement := fmt.Sprintf("INSERT INTO `lotto-info` (`draw-date`,`lotto-type`,`lotto-period`,`lotto-number-1`,`lotto-number-2`,`lotto-number-3`) VALUES('%s', '%s', '%s', '%s', '%s', '%s')", date, "3D", period, number1, number2, number3)
	_, err := db.Exec(statement)
	if err != nil {
		println(err)
		return
	}
	fmt.Println("today 3D lotto number inserted")
}
func lottoDiaryTask4D() {
	var date string
	var period string
	var number1 string
	var number2 string
	var number3 string
	var number4 string

	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	c := colly.NewCollector()
	c.OnHTML("div.contents_box04 > div#contents_logo_09", func(h *colly.HTMLElement) {
		date = strings.Split(h.DOM.Next().First().Text(), "\u00a0")[0]
		rex := regexp.MustCompile("[0-9]+")
		period = rex.FindString(strings.Split(h.DOM.Next().First().Text(), "\u00a0")[1])
		number := h.DOM.NextAll().Filter(".ball_tx").Text()
		number1 = number[0:1]
		number2 = number[1:2]
		number3 = number[2:3]
		number4 = number[3:4]
	})
	c.Visit(url)

	statement := fmt.Sprintf("INSERT INTO `lotto-info` (`draw-date`,`lotto-type`,`lotto-period`,`lotto-number-1`,`lotto-number-2`,`lotto-number-3`, `lotto-number-4`) VALUES('%s', '%s', '%s', '%s', '%s', '%s' , '%s')", date, "3D", period, number1, number2, number3, number4)
	_, err := db.Exec(statement)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println("today 4D lotto number inserted")
}

func DiaryTask() {

	task := gocron.NewScheduler()
	task.Every(1).Days().At("21:30:00").Do(lottoDiaryTask4D)
	task.Every(1).Days().At("21:30:00").Do(lottoDiaryTask3D)
	<-task.Start()
}

const (
	UserName     string = "root"
	Password     string = ""
	Addr         string = "db"
	Port         int    = 3306
	Database     string = "lotto"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

func getDatabaseConn() {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)

	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
	}
	db, err = conn.DB()
	if err != nil {
		fmt.Println("get db failed:", err)
	}
}

func insertLottoInfo() {

	sqlStatements, _ := os.ReadFile("data/db.sql")
	sqlStatementArray := strings.Split(string(sqlStatements), ";\n")

	for _, statement := range sqlStatementArray {
		_, err := db.Exec(statement)
		if err != nil {
			fmt.Println(err)
			db.Close()
			return
		}
	}
	fmt.Println("success migration")

	for year := 103; year < 112; year++ {
		for month := 1; month < 13; month++ {

			lotto3d := search3DNumber(year, month)
			for _, curLottoInfo := range lotto3d {
				drawdate_int, _ := strconv.ParseInt(curLottoInfo.drawDate, 10, 0)
				drawdate := time.Unix(int64(drawdate_int), 0).Format("2006-01-02 15:04:05")
				lottotype := curLottoInfo.lottoType
				period := curLottoInfo.period
				lottoNumber1 := curLottoInfo.lottoNumber[0:1]
				lottoNumber2 := curLottoInfo.lottoNumber[1:2]
				lottoNumber3 := curLottoInfo.lottoNumber[2:3]
				statement := fmt.Sprintf("INSERT INTO `lotto-info` (`draw-date`,`lotto-type`,`lotto-period`,`lotto-number-1`,`lotto-number-2`,`lotto-number-3`) VALUES('%s', '%s', '%s', '%s', '%s', '%s')", drawdate, lottotype, period, lottoNumber1, lottoNumber2, lottoNumber3)
				_, err := db.Exec(statement)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			lotto4d := search4DNumber(year, month)
			for _, curLottoInfo := range lotto4d {
				drawdate_int, _ := strconv.ParseInt(curLottoInfo.drawDate, 10, 0)
				drawdate := time.Unix(int64(drawdate_int), 0).Format("2006-01-02 15:04:05")
				lottotype := curLottoInfo.lottoType
				period := curLottoInfo.period
				lottoNumber1 := curLottoInfo.lottoNumber[0:1]
				lottoNumber2 := curLottoInfo.lottoNumber[1:2]
				lottoNumber3 := curLottoInfo.lottoNumber[2:3]
				lottoNumber4 := curLottoInfo.lottoNumber[3:4]
				statement := fmt.Sprintf("INSERT INTO `lotto-info` (`draw-date`,`lotto-type`,`lotto-period`,`lotto-number-1`,`lotto-number-2`,`lotto-number-3`,`lotto-number-4`) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s')", drawdate, lottotype, period, lottoNumber1, lottoNumber2, lottoNumber3, lottoNumber4)
				_, err := db.Exec(statement)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

		}
	}
	fmt.Println("success insert data")
	db.Close()

}
func main() {

	getDatabaseConn()
	insertLottoInfo()
	DiaryTask()
	db.Close()
}
