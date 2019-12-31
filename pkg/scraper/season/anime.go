package season

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/pkg/entities"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Anime struct {
	MalId       int                  `json:"mal_id"`
	Url         string               `json:"url"`
	Title       string               `json:"title"`
	ImageUrl    string               `json:"image_url"`
	Synopsis    string               `json:"synopsis"`
	Type        string               `json:"type"`
	AiringStart time.Time            `json:"airing_start"`
	Episodes    int                  `json:"episodes"`
	Members     int                  `json:"members"`
	Genres      []entities.MalEntity `json:"genres"`
	Source      string               `json:"source"`
	Producers   []entities.MalEntity `json:"producers"`
	Score       float64              `json:"score"`
	Licensors   []entities.MalEntity `json:"licensors"`
	R18         bool                 `json:"r18"`
	Kids        bool                 `json:"kids"`
	Continuing  bool                 `json:"continuing"`
}

type Season struct {
	SeasonName string   `json:"season_name"`
	SeasonYear string   `json:"season_year"`
	Anime      []*Anime `json:"anime"`
}

var scraper scrp.Scraper

func ScrapeAnime(year, season string) (*Season, error) {
	var url string
	if year == "" {
		url = "https://myanimelist.net/anime/season"
	} else {
		url = fmt.Sprintf("https://myanimelist.net/anime/season/%s/%s", year, season)
	}

	sel, err := scraper.GetSelection(url, ".seasonal-anime")
	if err != nil {
		return nil, err
	}

	var seas Season

	var wg sync.WaitGroup
	wg.Add(len(sel.Nodes))
	sel.Each(func(i int, s *goquery.Selection) {
		var anime Anime
		seas.Anime = append(seas.Anime, &anime)
		go func() {
			defer wg.Done()

			titleElem := s.Find(".link-title")
			anime.Title = titleElem.Text()
			anime.Url, _ = titleElem.Attr("href")

			anime.Source = s.Find(".source").Text()

			imgUrl, _ := s.Find(".image img").Attr("src")
			imgUrl = strings.ReplaceAll(imgUrl, "r/167x242/", "")
			imgUrl = strings.ReplaceAll(imgUrl, ".webp", ".jpg")
			anime.ImageUrl = strings.Split(imgUrl, "?")[0]

			anime.Synopsis = s.Find(".synopsis .preline").Text()

			anime.Score, _ = strconv.ParseFloat(strings.TrimSpace(s.Find(".score").Text()), 64)
			anime.Members, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(s.Find(".member").Text(), ",", "")))
		}()
	})

	wg.Wait()
	return &seas, nil
}
