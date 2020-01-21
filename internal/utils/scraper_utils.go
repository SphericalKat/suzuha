package utils

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func GetInfo(selection *goquery.Selection, info string) string {
	it := selection.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == info
	})
	parent := it.Parent()
	it.Remove()
	return strings.TrimSpace(parent.Text())
}

func GetInfoLinks(selection *goquery.Selection, info string) *goquery.Selection {
	it := selection.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == info
	})
	return it.SiblingsFiltered("a")
}

func GetInfoLinkString(selection *goquery.Selection, info string) string {
	it := selection.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == info
	})

	URL, _ := it.Attr("href")
	return URL
}