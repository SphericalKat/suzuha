package season

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"github.com/deletescape/suzuha/pkg/scraper/common"
	"sync"
)

type Season struct {
	SeasonName string       `json:"season_name"`
	SeasonYear string       `json:"season_year"`
	Anime      []*common.AnimeCard `json:"anime"`
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
		var anime common.AnimeCard
		seas.Anime = append(seas.Anime, &anime)
		go func() {
			defer wg.Done()
			common.ScrapeAnimeCard(&anime, s)
		}()
	})

	wg.Wait()
	return &seas, nil
}
