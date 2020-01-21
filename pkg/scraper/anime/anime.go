package anime

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/deletescape/suzuha/internal/utils"
	"github.com/deletescape/suzuha/pkg/entities"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)



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
		titleEnglish := utils.GetInfo(darkText, "English:")
		if titleEnglish != "" {
			anime.TitleEnglish = &titleEnglish
		}
		titleJapanese := utils.GetInfo(darkText, "Japanese:")
		if titleJapanese != "" {
			anime.TitleJapanese = &titleJapanese
		}
		synonyms := utils.GetInfo(darkText, "Synonyms:")
		if synonyms != "" {
			anime.TitleSynonyms = strings.Split(synonyms, ", ")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		atype := utils.GetInfo(darkText, "Type:")
		if atype != "" {
			anime.Type = &atype
		}
		episodes, err := strconv.Atoi(utils.GetInfo(darkText, "Episodes:"))
		if err == nil {
			anime.Episodes = &episodes
		} else {
			fmt.Println(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		rating := utils.GetInfo(darkText, "Rating:")
		if rating != "" {
			anime.Rating = &rating
		}
		duration := utils.GetInfo(darkText, "Duration:")
		if duration != "" {
			anime.Duration = &duration
		}
		source := utils.GetInfo(darkText, "Source:")
		if source != "" {
			anime.Source = &source
		}
		broadcast := utils.GetInfo(darkText, "Broadcast:")
		if broadcast != "" {
			anime.Broadcast = &broadcast
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		premiered := utils.GetInfo(darkText, "Premiered:")
		if premiered != "" {
			anime.Premiered = &premiered
		}
		status := utils.GetInfo(darkText, "Status:")
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
		popularity, err := strconv.Atoi(strings.TrimPrefix(utils.GetInfo(darkText, "Popularity:"), "#"))
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
		members, err := strconv.Atoi(strings.ReplaceAll(utils.GetInfo(darkText, "Members:"), ",", ""))
		if err == nil {
			anime.Members = &members
		} else {
			fmt.Println(err)
		}
		favorites, err := strconv.Atoi(strings.ReplaceAll(utils.GetInfo(darkText, "Favorites:"), ",", ""))
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
		utils.GetInfoLinks(darkText, "Studios:").Each(func(i int, s *goquery.Selection) {
			URL, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(URL)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(URL)
				if err == nil {
					anime.Studios = append(anime.Studios, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetInfoLinks(darkText, "Licensors:").Each(func(i int, s *goquery.Selection) {
			URL, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(URL)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(URL)
				if err == nil {
					anime.Licensors = append(anime.Licensors, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetInfoLinks(darkText, "Producers:").Each(func(i int, s *goquery.Selection) {
			URL, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(URL)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(URL)
				if err == nil {
					anime.Producers = append(anime.Producers, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetInfoLinks(darkText, "Genres:").Each(func(i int, s *goquery.Selection) {
			URL, _ := s.Attr("href")
			match := config.InfoLinkRe.FindStringSubmatch(URL)
			if len(match) == 3 {
				studioId, _ := strconv.Atoi(match[2])
				absolute, err := requestedUrl.Parse(URL)
				if err == nil {
					anime.Genres = append(anime.Genres, entities.MalEntity{MalId: studioId, Type: match[1], Name: s.Text(), Url: absolute.String()})
				}
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sel.Find(".theme-songs.opnening .theme-song").Each(func(i int, s *goquery.Selection) { // this is not a typo
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
				URL, _ := item.Attr("href")
				match := config.InfoLinkRe.FindStringSubmatch(URL)
				if len(match) == 3 {
					itemId, _ := strconv.Atoi(match[2])
					absolute, err := requestedUrl.Parse(URL)
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
			String: utils.GetInfo(darkText, "Aired:"),
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
