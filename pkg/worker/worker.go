package worker

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetLinks(doc Document) []string {
	body, err := getLinks(doc)
	if err != nil {
		return nil
	}
	return body
}

func getLinks(doc Document) ([]string, error) {
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
	for _, link := range links {
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
