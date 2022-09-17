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

func task() {
	fmt.Println("task")
}

func yearConvert(year string) (string, error) {
	ceYear, err := strconv.ParseInt(year, 10, 0)
	if err != nil {
		return "", err
	}
	ceYear = ceYear + 1911
	return strconv.Itoa(int(ceYear)), nil
}

func search3DNumber(twYear int, month int) {
	three_url := "https://www.taiwanlottery.com.tw/Lotto/3D/history.aspx"
	var formData = map[string]string{
		"L3DControl_history1$chk":       "radYM",
		"L3DControl_history1$dropYear":  strconv.Itoa(int(twYear)),
		"L3DControl_history1$dropMonth": strconv.Itoa(int(month)),
		"L3DControl_history1$btnSubmit": "查詢",
	}
	var lottoData = map[string]string{}
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
		twYear, err := yearConvert(date[0])
		if err != nil {
			fmt.Println(err)
		}
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
		lottoData[stempStr] = content

	})

	err := c.Post(three_url, formData)

	if err != nil {
		fmt.Println(err)
	}

	for period, number := range lottoData {
		fmt.Printf("%v: %v\n", period, number)
	}

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
	c.Visit(url)
	fmt.Println(date, number)
}

func diaryTask() {
	s := gocron.NewScheduler()
	s.Every(1).Days().At("21:00:00").Do(lottoDiaryTask3D)
	<-s.Start()
}

func test(taskCase int) {
	switch taskCase {
	case 0:
		s := gocron.NewScheduler()

		s.Every(1).Second().Do(lottoDiaryTask3D)
		<-s.Start()
	case 1:
		search3DNumber(111, 1)
	case 2:
		lottoDiaryTask3D()
	}
}

func main() {
	diaryTask()
	test(0)
}
