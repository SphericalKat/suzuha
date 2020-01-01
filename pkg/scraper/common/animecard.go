package common

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/deletescape/suzuha/pkg/entities"
	"strconv"
	"strings"
	"time"
)

type AnimeCard struct {
	MalId       int                  `json:"mal_id"`
	Url         string               `json:"url"`
	Title       string               `json:"title"`
	ImageUrl    string               `json:"image_url"`
	Synopsis    string               `json:"synopsis"`
	Type        string               `json:"type"`
	AiringStart time.Time            `json:"airing_start"`
	Episodes    *int                 `json:"episodes"`
	Members     int                  `json:"members"`
	Genres      []entities.MalEntity `json:"genres"`
	Source      string               `json:"source"`
	Producers   []entities.MalEntity `json:"producers"`
	Score       *float64             `json:"score"`
	Licensors   []string             `json:"licensors"`
	R18         bool                 `json:"r18"`
	Kids        bool                 `json:"kids"`
	Continuing  bool                 `json:"continuing"`
}

// 2006-01-02 15:04:05.999999999 -0700 MST
const timeFormat = "Jan 2, 2006, 15:04 (MST)"

func ScrapeAnimeCard(anime *AnimeCard, s *goquery.Selection) {
	titleElem := s.Find(".link-title")
	anime.Title = titleElem.Text()
	anime.Url, _ = titleElem.Attr("href")

	anime.Source = s.Find(".source").Text()

	imgUrl, _ := s.Find(".image img").Attr("src")
	imgUrl = strings.ReplaceAll(imgUrl, "r/167x242/", "")
	imgUrl = strings.ReplaceAll(imgUrl, ".webp", ".jpg")
	anime.ImageUrl = strings.Split(imgUrl, "?")[0]

	anime.Synopsis = s.Find(".synopsis .preline").Text()

	score, err := strconv.ParseFloat(strings.TrimSpace(s.Find(".score").Text()), 64)
	if err == nil {
		anime.Score = &score
	}
	anime.Members, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(s.Find(".member").Text(), ",", "")))

	id, _ := s.Find(".genres").Attr("id")
	anime.MalId, _ = strconv.Atoi(id)

	timeStr := strings.TrimSpace(s.Find(".remain-time").Text())
	anime.AiringStart, _ = time.Parse(timeFormat, timeStr)
	s.Find(".remain-time").Remove()

	anime.Type = strings.TrimSuffix(strings.TrimSpace(s.Find(".info").Text()), " -")

	episodes, err := strconv.Atoi(strings.TrimSuffix(s.Find(".eps span").Text(), " eps"))
	if err == nil {
		anime.Episodes = &episodes
	}

	anime.R18 = s.HasClass("r18")
	anime.Kids = s.HasClass("kids")
	anime.Continuing = strings.Contains(s.Parent().Find(".anime-header").Text(), "(Continuing)")

	lics, _ := s.Find(".licensors").Attr("data-licensors")
	licensors := strings.Split(lics, ",")
	for _, l := range licensors {
		if l != "" {
			anime.Licensors = append(anime.Licensors, l)
		}
	}

	genres := s.Find(".genre a")
	genres.Each(func(i int, sel *goquery.Selection) {

		u, _ := sel.Attr("href")
		absolute, err := config.MalUrl.Parse(u)
		if err == nil {
			match := config.InfoLinkRe.FindStringSubmatch(u)
			if len(match) == 3 {
				var genre entities.MalEntity

				genre.Type = match[1]
				genre.MalId, _ = strconv.Atoi(match[2])
				genre.Url = absolute.String()
				genre.Name = sel.Text()

				anime.Genres = append(anime.Genres, genre)
			}
		}
	})
}
