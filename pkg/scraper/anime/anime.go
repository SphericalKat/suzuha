package anime

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/toraberu/pkg/entities"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getInfo(selection *goquery.Selection, info string) string {
	it := selection.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == info
	})
	parent := it.Parent()
	it.Remove()
	return strings.TrimSpace(parent.Text())
}

func getInfoLinks(selection *goquery.Selection, info string) *goquery.Selection {
	it := selection.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == info
	})
	return it.SiblingsFiltered("a")
}

func cleanYtUrl(url string) string {
	match := ytLinkRe.FindStringSubmatch(url)
	if match != nil && len(match) > 1 {
		return fmt.Sprintf("https://youtu.be/%s", match[1])
	}
	return url
}

var ytLinkRe = regexp.MustCompile(`https?://(?:www\.)?youtube\.com/embed/([\w-]+).*`)
var infoLinkRe = regexp.MustCompile(`/(\w+)/(?:\w+/)?(\d+)/.*`)

func ScrapeAnime(id int) (entities.Anime, error) {
	anime := entities.Anime{MalId: &id}

	coll := colly.NewCollector()

	coll.OnHTML("div#contentWrapper", func(e *colly.HTMLElement) {
		darkText := e.DOM.Find(".dark_text")

		title := e.DOM.Find("span[itemprop=name]").First().Text()
		if title != "" {
			anime.Title = &title
		}
		titleEnglish := getInfo(darkText, "English:")
		if titleEnglish != "" {
			anime.TitleEnglish = &titleEnglish
		}
		titleJapanese := getInfo(darkText, "Japanese:")
		if titleJapanese != "" {
			anime.TitleJapanese = &titleJapanese
		}
		synonyms := getInfo(darkText, "Synonyms:")
		if synonyms != "" {
			anime.TitleSynonyms = strings.Split(synonyms, ", ")
		}
		atype := getInfo(darkText, "Type:")
		if atype != "" {
			anime.Type = &atype
		}
		episodes, err := strconv.Atoi(getInfo(darkText, "Episodes:"))
		if err == nil {
			anime.Episodes = &episodes
		} else {
			fmt.Println(err)
		}
		rating := getInfo(darkText, "Rating:")
		if rating != "" {
			anime.Rating = &rating
		}
		duration := getInfo(darkText, "Duration:")
		if duration != "" {
			anime.Duration = &duration
		}
		source := getInfo(darkText, "Source:")
		if source != "" {
			anime.Source = &source
		}
		broadcast := getInfo(darkText, "Broadcast:")
		if broadcast != "" {
			anime.Broadcast = &broadcast
		}
		premiered := getInfo(darkText, "Premiered:")
		if premiered != "" {
			anime.Premiered = &premiered
		}
		status := getInfo(darkText, "Status:")
		if status != "" {
			anime.Status = &status
		}
		score, err := strconv.ParseFloat(e.ChildText("div[itemprop=aggregateRating] span[itemprop=ratingValue]"), 64)
		if err == nil {
			anime.Score = &score
		} else {
			fmt.Println(err)
		}
		scoredBy, err := strconv.Atoi(e.ChildText("div[itemprop=aggregateRating] span[itemprop=ratingCount]"))
		if err == nil {
			anime.ScoredBy = &scoredBy
		} else {
			fmt.Println(err)
		}
		popularity, err := strconv.Atoi(strings.TrimPrefix(getInfo(darkText, "Popularity:"), "#"))
		if err == nil {
			anime.Popularity = &popularity
		} else {
			fmt.Println(err)
		}
		rank, err := strconv.Atoi(strings.TrimPrefix(e.ChildText(".numbers.ranked strong"), "#"))
		if err == nil {
			anime.Rank = &rank
		} else {
			fmt.Println(err)
		}
		members, err := strconv.Atoi(strings.ReplaceAll(getInfo(darkText, "Members:"), ",", ""))
		if err == nil {
			anime.Members = &members
		} else {
			fmt.Println(err)
		}
		favorites, err := strconv.Atoi(strings.ReplaceAll(getInfo(darkText, "Favorites:"), ",", ""))
		if err == nil {
			anime.Favorites = &favorites
		} else {
			fmt.Println(err)
		}
		synopsys := e.ChildText("span[itemprop=description]")
		if synopsys != "" {
			anime.Synopsis = &synopsys
		}
		imageUrl := e.ChildAttr("img[itemprop=image]", "src")
		if imageUrl != "" {
			anime.ImageUrl = &imageUrl
		}
		url := e.ChildAttr("a.horiznav_active", "href")
		if url == "" {
			url = e.Request.URL.String()
		}
		anime.Url = &url
		trailerUrl := cleanYtUrl(e.ChildAttr("div.video-promotion a.video-unit", "href"))
		if trailerUrl != "" {
			anime.TrailerUrl = &trailerUrl
		}
		getInfoLinks(darkText, "Studios:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := infoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])

				anime.Studios = append(anime.Studios, entities.Studio{MalId: studioId, Type: match[1], Name: s.Text(), Url: e.Request.AbsoluteURL(url)})
			}
		})
		getInfoLinks(darkText, "Licensors:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := infoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])

				anime.Licensors = append(anime.Licensors, entities.Studio{MalId: studioId, Type: match[1], Name: s.Text(), Url: e.Request.AbsoluteURL(url)})
			}
		})
		getInfoLinks(darkText, "Producers:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := infoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])

				anime.Producers = append(anime.Producers, entities.Studio{MalId: studioId, Type: match[1], Name: s.Text(), Url: e.Request.AbsoluteURL(url)})
			}
		})
		getInfoLinks(darkText, "Genres:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := infoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])

				anime.Genres = append(anime.Genres, entities.Genre{MalId: studioId, Type: match[1], Name: s.Text(), Url: e.Request.AbsoluteURL(url)})
			}
		})
		e.DOM.Find(".theme-songs.opnening .theme-song").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.HasPrefix(text, "#") {
				text = strings.SplitAfterN(text, ": ", 2)[1]
			}
			anime.OpeningThemes = append(anime.OpeningThemes, text)
		})
		e.DOM.Find(".theme-songs.ending .theme-song").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.HasPrefix(text, "#") {
				text = strings.SplitAfterN(text, ": ", 2)[1]
			}
			anime.EndingThemes = append(anime.EndingThemes, text)
		})
		anime.Related = map[string][]entities.RelatedItem{}
		e.DOM.Find(".anime_detail_related_anime tr").Each(func(i int, s *goquery.Selection) {
			key := strings.TrimSuffix(s.Find("td").First().Text(), ":")
			s.Find("td a").Each(func(i int, item *goquery.Selection) {
				url, _ := item.Attr("href")
				match := infoLinkRe.FindStringSubmatch(url)
				if len(match) == 3 {
					itemId, _ := strconv.Atoi(match[2])

					anime.Related[key] = append(anime.Related[key], entities.RelatedItem{
						MalId: itemId,
						Type:  match[1],
						Name:  item.Text(),
						Url:   e.Request.AbsoluteURL(url),
					})
				}
			})
		})

		// Clean up to get background text
		backgroundTitle := e.DOM.Find("td h2").FilterFunction(func(i int, s *goquery.Selection) bool {
			return strings.Contains(s.Text(), "Background")
		})
		backgroundParent := backgroundTitle.Parent()
		backgroundParent.Children().First().NextUntilSelection(backgroundTitle).Remove()
		backgroundParent.Children().First().Remove()
		backgroundTitle.Remove()
		backgroundParent.RemoveFiltered("div.border_top")
		background := strings.TrimSpace(backgroundParent.Text())
		if background != "" && background != "No background information has been added to this title. Help improve our database by adding background information here." {
			anime.Background = &background
		}

		anime.Airing = *anime.Status == "Currently Airing"

		// TODO: HUGE CRIMES AHEAD HOW TF DO I CLEAN THIS UP
		anime.Aired = entities.Aired{
			String: getInfo(darkText, "Aired:"),
			Prop: entities.AiredProp{
				From: entities.Date{},
				To:   entities.Date{},
			},
		}
		airedParts := strings.Split(anime.Aired.String, " to ")
		tmpFrom, err := time.Parse("Jan _2, 2006", airedParts[0])
		if err == nil {
			anime.Aired.From = &tmpFrom
			tmpDay := anime.Aired.From.Day()
			anime.Aired.Prop.From.Day = &tmpDay
			tmpMonth := int(anime.Aired.From.Month())
			anime.Aired.Prop.From.Month = &tmpMonth
			tmpYear := anime.Aired.From.Year()
			anime.Aired.Prop.From.Year = &tmpYear
		} else {
			fmt.Println(err)
		}
		if len(airedParts) == 2 {
			tmpTo, err := time.Parse("Jan _2, 2006", airedParts[1])
			if err == nil {
				anime.Aired.To = &tmpTo
				tmpDay := anime.Aired.To.Day()
				anime.Aired.Prop.To.Day = &tmpDay
				tmpMonth := int(anime.Aired.To.Month())
				anime.Aired.Prop.To.Month = &tmpMonth
				tmpYear := anime.Aired.To.Year()
				anime.Aired.Prop.To.Year = &tmpYear
			} else {
				fmt.Println(err)
			}
		}
	})

	err := coll.Visit(fmt.Sprintf("https://myanimelist.net/anime/%d", id))
	return anime, err
}
