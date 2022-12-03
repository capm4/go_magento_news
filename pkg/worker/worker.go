package worker

import (
	"fmt"
	"magento/bot/pkg/model"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

func GetLinks(doc model.Website) []string {
	body, err := getLinks(doc)
	if err != nil {
		return nil
	}
	return body
}

func getLinks(doc model.Website) ([]string, error) {
	var links []string
	index := 0
	response, err := http.Get(doc.Url)
	if err != nil {
		logrus.Warning(fmt.Sprintf("error while get link. Error : %s", err))
		return []string{}, err
	}
	docReader, err := goquery.NewDocumentFromReader(response.Body)
	docReader.Find(doc.Selector).Each(func(i int, s *goquery.Selection) {
		// For each item found, get href value
		href, _ := s.Attr(doc.Attribute)
		if href != "" {
			links = append(links, href)
		}
	})
	uniqueLinks := unique(links)
	for _, link := range uniqueLinks {
		if link == doc.LastUrl {
			break
		}
		index++
	}
	if index == 0 {
		return []string{}, nil
	}
	return links[:index], nil
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
