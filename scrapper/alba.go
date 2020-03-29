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
	Address string
	IsImg   bool
}

type Data struct {
	TotalItem int
	Result []Alba
}

// GetAlbaPages scrap http://alba.co.kr search result
func GetAlbaPages(job string, area string, volume int) (*Data, error) {
	jobQuery := url.QueryEscape(convertUTF8ToEUCKR(job))
	areaQuery := url.QueryEscape(convertUTF8ToEUCKR(area))
	search := fmt.Sprintf("/search/Search.asp?WsSrchWord=%s&wsSrchWordarea=%s&Section=0&Page=1&hidSort=FREEORDER&hidSortOrder=1&hidSortDate=rday0&hidSortCnt=%d&gendercd=C03&ageconst=G01", jobQuery, areaQuery, volume)
	var result []Alba
	c := make(chan []Alba)

	res, err := http.Get(baseURL + search)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	totalItemText := doc.Find("#SearchJob>h2>em").Text()
	totalItemInt, err := strconv.Atoi(strings.Replace(totalItemText, ",", "", -1))
	if err != nil {
		return nil, err
	}

	totalPages := totalItemInt / volume
	if totalItemInt%volume > 0 {
		totalPages++
	}

	for i := 1; i <= totalPages; i++ {
		searchEachPages := fmt.Sprintf("/search/Search.asp?WsSrchWord=%s&wsSrchWordarea=%s&Section=0&Page=%d&hidSort=FREEORDER&hidSortOrder=1&hidSortDate=rday0&hidSortCnt=%d&gendercd=C03&ageconst=G01", jobQuery, areaQuery, i, volume)
		go getPages(searchEachPages, c)
	}

	for i := 1; i <=totalPages; i++ {
		itemList := <-c
		result = append(result, itemList...)
	}

	data := Data{
		TotalItem: len(result),
		Result: result,
	}

	return &data, nil
}

func getPages(url string, mainC chan<- []Alba) error {
	c := make(chan Alba)
	var result []Alba

	res, err := http.Get(baseURL + url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	jobResult := doc.Find("#jobNormal")
	jobResult.Find("li").Each(func(i int, li *goquery.Selection) {
		go extractJob(li, c)
	})

	for i := 0; i < jobResult.Find("li").Length(); i++ {
		alba := <-c
		result = append(result, alba)
	}

	mainC <- result
	return nil
}

func extractJob(li *goquery.Selection, c chan<- Alba) {
	title := li.Find(".title>a").Text()
	link, _ := li.Find(".title>a").Attr("href")

	splitLink := strings.Split(link, "=")
	adid, _ := strconv.Atoi(splitLink[1])
	address, isImg, _ := getAlbaAddress(baseURL + link)

	c <- Alba{
		Title: convertEUCKRToUTF8(title),
		Adid:  adid,
		Address: address,
		IsImg: isImg,
	}
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
