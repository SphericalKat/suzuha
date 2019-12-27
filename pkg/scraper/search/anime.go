package search

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	scrp "github.com/deletescape/toraberu/pkg/scraper"
	"strconv"
	"strings"
)

type Animes struct {
	LastPage int     `json:"last_page"`
	Results  []Anime `json:"results"`
}

type Anime struct {
	MalId    int     `json:"mal_id"`
	Url      string  `json:"url"`
	ImageUrl string  `json:"image_url"`
	Title    string  `json:"title"`
	Airing   bool    `json:"airing"`
	Synopsis string  `json:"synopsis"`
	Type     string  `json:"type"`
	Episodes int     `json:"episodes"`
	Score    float64 `json:"score"`
	//StartDate time.Time `json:"start_date"`
	//EndDate   time.Time `json:"end_date"`
	Members int    `json:"members"`
	Rated   string `json:"rated"`
}

var scraper scrp.Scraper

func ScrapeAnimeSearch(query string, page int) (*Animes, error) {
	searchUrl := fmt.Sprintf("https://myanimelist.net/anime.php?q=%s&show=%d", query, (page-1)*50)
	sel, err := scraper.GetSelection(searchUrl, ".js-block-list tr")
	if err != nil {
		return nil, err
	}
	var animes Animes
	sel.Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}
		var anime Anime
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			switch i {
			case 0:
				link := td.Find(".hoverinfo_trigger")
				anime.Url, _ = link.Attr("href")
				// Get the url of the original image and clean it up
				imgUrl, _ := td.Find("img").Attr("data-src")
				imgUrl = strings.ReplaceAll(imgUrl, "r/50x70/", "")
				imgUrl = strings.ReplaceAll(imgUrl, "r/100x140/", "")
				imgUrl = strings.ReplaceAll(imgUrl, ".webp", ".jpg")
				anime.ImageUrl = strings.Split(imgUrl, "?")[0]
				break
			case 1:
				titleElem := td.Find(".hoverinfo_trigger")
				anime.Title = titleElem.Text()
				anime.MalId, _ = strconv.Atoi(strings.TrimPrefix(titleElem.AttrOr("id", ""), "sinfo"))
				break
			case 2:
				anime.Type = strings.TrimSpace(td.Text())
				break
			case 3:
				anime.Episodes, _ = strconv.Atoi(strings.TrimSpace(td.Text()))
				break
			case 4:
				anime.Score, _ = strconv.ParseFloat(strings.TrimSpace(td.Text()), 64)
				break
			}
		})
		animes.Results = append(animes.Results, anime)
	})

	return &animes, err
}
