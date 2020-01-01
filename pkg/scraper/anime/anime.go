package anime

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/deletescape/suzuha/pkg/entities"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

var scraper scrp.Scraper
var ytLinkRe = regexp.MustCompile(`https?://(?:www\.)?youtube\.com/embed/([\w-]+).*`)

func ScrapeAnime(id int) (*entities.Anime, error) {
	anime := entities.Anime{MalId: &id}

	requested := fmt.Sprintf("https://myanimelist.net/anime/%d", id)
	sel, err := scraper.GetSelection(requested, "div#contentWrapper")
	if err != nil {
		return nil, err
	}
	requestedUrl, err := url.Parse(requested)
	if err != nil {
		return nil, err
	}

	darkText := sel.Find(".dark_text")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		title := sel.Find("span[itemprop=name]").First().Text()
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
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		premiered := getInfo(darkText, "Premiered:")
		if premiered != "" {
			anime.Premiered = &premiered
		}
		status := getInfo(darkText, "Status:")
		if status != "" {
			anime.Status = &status
			anime.Airing = status == "Currently Airing"
		}
		synopsys := sel.Find("span[itemprop=description]").Text()
		if synopsys != "" {
			anime.Synopsis = &synopsys
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scoredBy, err := strconv.Atoi(sel.Find("div[itemprop=aggregateRating] span[itemprop=ratingCount]").Text())
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
		rank, err := strconv.Atoi(strings.TrimPrefix(sel.Find(".numbers.ranked strong").Text(), "#"))
		if err == nil {
			anime.Rank = &rank
		} else {
			fmt.Println(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
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
		score, err := strconv.ParseFloat(sel.Find("div[itemprop=aggregateRating] span[itemprop=ratingValue]").Text(), 64)
		if err == nil {
			anime.Score = &score
		} else {
			fmt.Println(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		imageUrl, ok := sel.Find("img[itemprop=image]").Attr("src")
		if ok || imageUrl != "" {
			anime.ImageUrl = &imageUrl
		}
		animeUrl, ok := sel.Find("a.horiznav_active").Attr("href")
		if !ok || animeUrl == "" {
			animeUrl = requested
		}
		anime.Url = &animeUrl
		trailerUrl, ok := sel.Find("div.video-promotion a.video-unit").Attr("href")
		if ok && trailerUrl != "" {
			trailerUrl = cleanYtUrl(trailerUrl)
			anime.TrailerUrl = &trailerUrl
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		getInfoLinks(darkText, "Studios:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(url)
				if err == nil {
					anime.Studios = append(anime.Studios, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		getInfoLinks(darkText, "Licensors:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(url)
				if err == nil {
					anime.Licensors = append(anime.Licensors, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		getInfoLinks(darkText, "Producers:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(url)
				if err == nil {
					anime.Producers = append(anime.Producers, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		getInfoLinks(darkText, "Genres:").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(url)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(url)
				if err == nil {
					anime.Genres = append(anime.Genres, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sel.Find(".theme-songs.opnening .theme-song").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.HasPrefix(text, "#") {
				text = strings.SplitAfterN(text, ": ", 2)[1]
			}
			anime.OpeningThemes = append(anime.OpeningThemes, text)
		})
		sel.Find(".theme-songs.ending .theme-song").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.HasPrefix(text, "#") {
				text = strings.SplitAfterN(text, ": ", 2)[1]
			}
			anime.EndingThemes = append(anime.EndingThemes, text)
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		anime.Related = map[string][]entities.MalEntity{}
		sel.Find(".anime_detail_related_anime tr").Each(func(i int, s *goquery.Selection) {
			key := strings.TrimSuffix(s.Find("td").First().Text(), ":")
			s.Find("td a").Each(func(i int, item *goquery.Selection) {
				url, _ := item.Attr("href")
				match := config.InfoLinkRe.FindStringSubmatch(url)
				if len(match) == 3 {
					itemId, _ := strconv.Atoi(match[2])
					absolute, err := requestedUrl.Parse(url)
					if err == nil {
						anime.Related[key] = append(anime.Related[key], entities.MalEntity{
							MalId: itemId,
							Type:  match[1],
							Name:  item.Text(),
							Url:   absolute.String(),
						})
					}
				}
			})
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Clean up to get background text
		backgroundTitle := sel.Find("td h2").FilterFunction(func(i int, s *goquery.Selection) bool {
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
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		anime.Aired = entities.Aired{
			String: getInfo(darkText, "Aired:"),
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
	}()

	wg.Wait()
	return &anime, err
}
