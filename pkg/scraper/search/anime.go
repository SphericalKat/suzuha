package search

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	scrp "github.com/deletescape/toraberu/pkg/scraper"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Animes struct {
	LastPage int      `json:"last_page"`
	Results  []*Anime `json:"results"`
}

type Anime struct {
	MalId     int        `json:"mal_id"`
	Url       string     `json:"url"`
	ImageUrl  string     `json:"image_url"`
	Title     string     `json:"title"`
	Airing    bool       `json:"airing"`
	Synopsis  string     `json:"synopsis"`
	Type      string     `json:"type"`
	Episodes  int        `json:"episodes"`
	Score     float64    `json:"score"`
	Members   int        `json:"members"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Rated     string     `json:"rated"`
}

var scraper scrp.Scraper

func ScrapeAnimeSearch(query string, page int) (*Animes, error) {
	searchUrl := fmt.Sprintf("https://myanimelist.net/anime.php?q=%s&show=%d&c[]=a&c[]=b&c[]=c&c[]=d&c[]=e&c[]=f&c[]=g", query, (page-1)*50)
	sel, err := scraper.GetSelection(searchUrl, "#content")
	if err != nil {
		return nil, err
	}

	trs := sel.Find(".js-block-list tr")
	count := len(trs.Nodes) - 1
	animes := Animes{
		Results: make([]*Anime, count),
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		lastPage := sel.Find(".normal_header.pt16 span a:last-of-type").Text()
		if lastPage == "" {
			animes.LastPage = 1
		} else {
			animes.LastPage, _ = strconv.Atoi(lastPage)
		}
	}()

	wg.Add(count)
	trs.Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}
		go func() {
			defer wg.Done()
			var anime Anime
			tr.Find("td").Each(func(i int, td *goquery.Selection) {
				switch i {
				case 0:
					link := td.Find(".hoverinfo_trigger")
					anime.Url, _ = link.Attr("href")
					anime.MalId, _ = strconv.Atoi(strings.TrimPrefix(link.AttrOr("id", ""), "sarea"))
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
					synopsis := td.Find("div")
					synopsis.Find("a").Remove()
					anime.Synopsis = strings.TrimSpace(synopsis.Text())
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
				case 5:
					dateStr := strings.TrimSpace(td.Text())
					dateStr = strings.ReplaceAll(dateStr, "??", "01")
					dateStr = strings.ReplaceAll(dateStr, "?", "0")
					date, err := time.Parse("01-02-06", dateStr)
					if err == nil {
						anime.StartDate = &date
					}
					break
				case 6:
					dateStr := strings.TrimSpace(td.Text())
					dateStr = strings.ReplaceAll(dateStr, "??", "01")
					dateStr = strings.ReplaceAll(dateStr, "?", "0")
					date, err := time.Parse("01-02-06", dateStr)
					now := time.Now()
					if err == nil {
						anime.EndDate = &date
						anime.Airing = date.After(now) && anime.StartDate.Before(now)
					} else {
						anime.Airing = anime.StartDate == nil || anime.StartDate.Before(now)
					}
					break
				case 7:
					anime.Members, _ = strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(td.Text()), ",", ""))
					break
				case 8:
					anime.Rated = strings.TrimSpace(td.Text())
					break
				}
			})

			animes.Results[i-1] = &anime
		}()
	})
	wg.Wait()
	return &animes, err
}
