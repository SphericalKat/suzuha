package search

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type Animes struct {
	LastPage int     `json:"last_page"`
	Results  []Anime `json:"results"`
}

type Anime struct {
	MalId     int       `json:"mal_id"`
	Url       string    `json:"url"`
	ImageUrl  string    `json:"image_url"`
	Title     string    `json:"title"`
	Airing    bool      `json:"airing"`
	Synopsis  string    `json:"synopsis"`
	Type      string    `json:"type"`
	Episodes  int       `json:"episodes"`
	Score     float64   `json:"score"`
	//StartDate time.Time `json:"start_date"`
	//EndDate   time.Time `json:"end_date"`
	Members   int       `json:"members"`
	Rated     string    `json:"rated"`
}

func ScrapeAnimeSearch(query string, page int) (Animes, error) {
	var animes Animes

	coll := colly.NewCollector()
	coll.OnHTML(".js-block-list tr", func(e *colly.HTMLElement) {
		var anime Anime
		var isResult = false
		e.ForEachWithBreak("td", func(i int, td *colly.HTMLElement) bool {
			switch i {
			case 0:
				link := td.DOM.Find(".hoverinfo_trigger")
				url, exists := link.Attr("href")
				if !exists || url == "" {
					return false
				}
				isResult = true
				anime.Url = url
				// Get the url of the original image and clean it up
				imgUrl, _ := td.DOM.Find("img").Attr("data-src")
				imgUrl = strings.ReplaceAll(imgUrl, "r/50x70/", "")
				imgUrl = strings.ReplaceAll(imgUrl, "r/100x140/", "")
				imgUrl = strings.ReplaceAll(imgUrl, ".webp", ".jpg")
				anime.ImageUrl = strings.Split(imgUrl, "?")[0]
				break
			case 1:
				titleElem := td.DOM.Find(".hoverinfo_trigger")
				anime.Title = titleElem.Text()
				anime.MalId, _ = strconv.Atoi(strings.TrimPrefix(titleElem.AttrOr("id", ""), "sinfo"))
				break
			case 2:
				anime.Type = strings.TrimSpace(td.Text)
				break
			case 3:
				anime.Episodes, _ = strconv.Atoi(strings.TrimSpace(td.Text))
				break
			case 4:
				anime.Score, _ = strconv.ParseFloat(strings.TrimSpace(td.Text), 64)
				break
			}
			return true
		})
		if isResult {
			// TODO: we need to somehow
			/*hoverInfoColl := colly.NewCollector()
			hoverInfoColl.OnScraped(func(r *colly.Response) {
				doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
				synopsisElem := doc.Find(".hoverinfo-contaniner div")
				synopsisElem.Children().RemoveFiltered("a")
				anime.Synopsis = synopsisElem.Text()
				anime.Members = 5000
			})
			hoverInfoColl.Visit(fmt.Sprintf("https://myanimelist.net/includes/ajax.inc.php?t=64&id=%d", anime.MalId))*/
			animes.Results = append(animes.Results, anime)
		}
	})

	err := coll.Visit(fmt.Sprintf("https://myanimelist.net/anime.php?q=%s&show=%d", query, (page-1)*50))
	return animes, err
}
