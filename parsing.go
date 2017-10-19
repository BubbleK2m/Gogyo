package gogyo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const MonthlyMenuPath = "sts_sci_md00_001.do"

func GetDailyMenu(school School, year int, month int, date int) Menu {
	return GetMonthlyMenus(school, year, month)[date-1]
}

func GetMonthlyMenus(school School, year int, month int) map[int]Menu {
	return ParseMonthlyMenus(school, year, month)
}

func ParseMonthlyMenus(school School, year int, month int) map[int]Menu {
	monthly := make(map[int]Menu)

	path := fmt.Sprintf("http://%s/%s?schulCode=%s&schulCrseScCode=%d&schulKndScCode=0%d&schYm=%d%02d", school.Region, MonthlyMenuPath, school.Code, school.Kind, school.Kind, year, month)
	document, exception := goquery.NewDocument(path)

	if exception != nil {
		log.Fatal(exception)
	}

	selection := document.Find(".tbl_calendar tbody tr td div") // 급식 정보를 담는 요소 선택
	selection = selection.FilterFunction(func(index int, selection *goquery.Selection) bool {
		html, exception := selection.Html()

		if exception != nil {
			log.Fatal(exception)
		}

		html = strings.Trim(html, " ")
		chunks := strings.SplitAfter(html, "<br/>")

		return len(chunks) > 1
	})

	selection.Each(func(index int, selection *goquery.Selection) {
		html, exception := selection.Html()

		if exception != nil {
			log.Fatal(exception)
		}

		html = strings.Trim(html, " ")
		chunks := strings.Split(html, "<br/>")

		date, exception := strconv.Atoi(chunks[0])

		if exception != nil {
			log.Fatal(exception)
		}

		daily := make(Menu)

		times := [3]string{"[조식]", "[중식]", "[석식]"}
		var position int

	chunks:
		for _, chunk := range chunks[1:] {
			chunk = strings.Replace(chunk, "&amp", "&", -1)

			for index, time := range times {
				if strings.Compare(chunk, time) == 0 {
					position = index
					daily[position] = make([]Meal, 0)

					continue chunks
				}
			}

			chunks := strings.FieldsFunc(chunk, func(code rune) bool {
				return unicode.IsNumber(code) || !unicode.IsLetter(code)
			})

			for _, chunk := range chunks {
				meal := NewMeal(MealName(chunk))
				daily[position] = append(daily[position], *meal)
			}
		}

		monthly[date-1] = daily
	})

	return monthly
}
