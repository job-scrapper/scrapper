package scrapper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

var baseURL = "http://alba.co.kr"

// Alba is struct of http://alba.co.kr values
type Alba struct {
	Adid    int
	Title   string
	TelText string
	TelLink string
}

// GetAlbaPages scrap http://alba.co.kr search result
func GetAlbaPages(job string, area string) ([]Alba, error) {
	jobQuery := url.QueryEscape(convertUTF8ToEUCKR(job))
	areaQuery := url.QueryEscape(convertUTF8ToEUCKR(area))
	search := fmt.Sprintf("/search/Search.asp?WsSrchWord=%s&wsSrchWordarea=%s&Section=0&Page=1&hidSort=FREEORDER&hidSortOrder=1&hidSortDate=rday0&hidSortCnt=50&gendercd=C03&ageconst=G01", jobQuery, areaQuery)
	var result []Alba

	res, err := http.Get(baseURL + search)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	jobResult := doc.Find("#jobNormal")
	jobResult.Find("li").Each(func(i int, li *goquery.Selection) {
		albaItem := new(Alba)

		title := li.Find(".title>a").Text()
		link, _ := li.Find(".title>a").Attr("href")

		splitLink := strings.Split(link, "=")
		adid, _ := strconv.Atoi(splitLink[1])

		albaItem.Title = convertEUCKRToUTF8(title)
		albaItem.Adid = adid

		address, isImg, _ := getAlbaAddress(baseURL + link)
		
		if isImg {
			albaItem.TelLink = address
		} else {
			albaItem.TelText = address
		}

		result = append(result, *albaItem)
	})

	return result, nil
}

func getAlbaAddress(url string) (string, bool, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", false, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", false, err
	}

	telEmail := doc.Find("#InfoApply>.info>.telEmail")
	textTel := telEmail.Find(".contact").Text()
	link, isExist := telEmail.Find("img").Attr("src")

	if isExist {
		return "http:" + link, true, nil
	}
	return convertEUCKRToUTF8(textTel), false, nil
}

func convertEUCKRToUTF8(str string) string {
	out, _ := iconv.ConvertString(str, "EUC-KR", "UTF-8")
	return out
}

func convertUTF8ToEUCKR(str string) string {
	out, _ := iconv.ConvertString(str, "UTF-8", "EUC-KR")
	return out
}
