package scrapper

import (
	"net/http"
	"strings"
	"strconv"

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
func GetAlbaPages() ([]Alba, error) {
	search := "/search/Search.asp?WsSrchWord=%C1%D6%B9%E6&wsSrchWordarea=%BC%AD%BF%EF&Section=0&Page=&hidschContainText=&hidWsearchInOut=&hidGroupKeyJobArea=&hidGroupKeyJobHotplace=&hidGroupKeyJobJobKind=&hidGroupKeyResumeArea=&hidGroupKeyResumeJobKind=&hidGroupKeyPay=&hidGroupKeyWorkWeek=&hidGroupKeyWorkPeriod=&hidGroupKeyOpt=&hidGroupKeyGender=&hidGroupKeyAge=&hidGroupKeyCareer=&hidGroupKeyLicense=&hidGroupKeyEduData=&hidGroupKeyWorkTime=&hidGroupKeyWorkState=&hidGroupKeyJobCareer=&hidSort=&hidSortOrder=1&hidSortDate=&hidSortCnt=&hidSortFilter=&hidArea=&area=&hidJobKind=&jobkind=&gendercd=C03&ageconst=G01&agelimitmin=&agelimitmax=&workperiod=&workweek="
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
	} else {
		return convertEUCKRToUTF8(textTel), false, nil
	}
}

func convertEUCKRToUTF8(str string) string {
	out, _ := iconv.ConvertString(str, "EUC-KR", "UTF-8")
	return out
}
