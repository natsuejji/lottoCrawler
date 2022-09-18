package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/jasonlvhit/gocron"
)

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
func lottoDiaryTask3D() {
	var date string
	var number string
	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	c := colly.NewCollector()
	c.OnHTML("div.contents_box04 > div#contents_logo_08", func(h *colly.HTMLElement) {
		date = strings.Split(h.DOM.Next().First().Text(), "\u00a0")[0]
		number = h.DOM.NextAll().Filter(".ball_tx").Text()
	})
	fmt.Println(date, number)
	c.Visit(url)
}
func search4DNumber(twYear int, month int) []LottoInfo {
	three_url := "https://www.taiwanlottery.com.tw/Lotto/4D/history.aspx"
	var formData = map[string]string{
		"L3DControl_history1$chk":       "radYM",
		"L3DControl_history1$dropYear":  strconv.Itoa(int(twYear)),
		"L3DControl_history1$dropMonth": strconv.Itoa(int(month)),
		"L3DControl_history1$btnSubmit": "查詢",
	}
	var lottoDatas []LottoInfo
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

	err := c.Post(three_url, formData)

	if err != nil {
		fmt.Println(err)
	}

	return lottoDatas
}
func lottoDiaryTask4D() {
	var date string
	var number string
	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	c := colly.NewCollector()
	c.OnHTML("div.contents_box04 > div#contents_logo_09", func(h *colly.HTMLElement) {
		date = strings.Split(h.DOM.Next().First().Text(), "\u00a0")[0]
		number = h.DOM.NextAll().Filter(".ball_tx").Text()
	})
	fmt.Println(date, number)
	c.Visit(url)
}

func DiaryTask() {

	task := gocron.NewScheduler()
	task.Every(1).Days().At("21:01:00").Do(lottoDiaryTask4D)
	task.Every(1).Days().At("21:01:00").Do(lottoDiaryTask3D)
	<-task.Start()
}

func main() {

	// lotto3d := search3DNumber(111, 9)
	// lotto4d := search4DNumber(111, 9)

	DiaryTask()
}
