package person

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/internal/utils"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
)

type Person struct {
	URL            *string    `json:"url"`
	ImageURL       *string    `json:"image_url"`
	WebsiteURL     *string    `json:"website_url"`
	GivenName      *string    `json:"given_name"`
	FamilyName     *string    `json:"family_name"`
	About          *string    `json:"about"`
	Name           *string    `json:"name"`
	MalID          *int       `json:"mal_id"`
	MemberFaves    *int       `json:"member_favorites"`
	AlternateNames []string   `json:"alternate_names"`
	Birthday       *time.Time `json:"birthday"`
	Roles          []*VARole  `json:"voice_acting_roles"`
}

type VARole struct {
	Role  string   `json:"role"`
	Anime *VAAnime `json:"anime"`
	Char  *VAChar  `json:"character"`
}

type VAAnime struct {
	MalID    *int    `json:"mal_id"`
	URL      *string `json:"url"`
	ImageURL *string `json:"image_url"`
}

type VAChar struct {
	MalID    *int    `json:"mal_id"`
	URL      *string `json:"url"`
	ImageURL *string `json:"image_url"`
	Name     *string `json:"name"`
}

var scraper scrp.Scraper

const layoutUS = "Jan  2, 2006"

func ScrapePerson(id int) (*Person, error) {
	person := &Person{MalID: &id}
	requested := fmt.Sprintf("https://myanimelist.net/people/%d", id)
	sel, err := scraper.GetSelection(requested, "div#contentWrapper")
	if err != nil {
		return nil, err
	}
	doc, err := scraper.Get(requested)
	if err != nil {
		return nil, err
	}
	_, err = url.Parse(requested)
	if err != nil {
		return nil, err
	}

	darkText := sel.Find(".dark_text")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		navButtons := sel.Find(".horiznav_active")
		imageURL, ok := sel.Find("td.borderClass").Find("img.lazyload").Attr("data-src")

		URL := utils.GetInfoLinkString(navButtons, "Details")
		if ok || imageURL != "" {
			person.ImageURL = &imageURL
		}
		if URL != "" {
			person.URL = &URL
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := utils.GetInfoLinks(darkText, "Website:")
		websiteURL, _ := s.Attr("href")
		givenName := utils.GetInfo(darkText, "Given name:")
		memberFaves, err := strconv.Atoi(strings.ReplaceAll(utils.GetInfo(darkText, "Member Favorites:"), ",", ""))
		altNames := strings.Split(utils.GetInfo(darkText, "Alternate names:"), ", ")

		bday := utils.GetInfo(darkText, "Birthday:")
		bdayT, errT := time.Parse(layoutUS, bday)

		if websiteURL != "" && websiteURL != "http://" && websiteURL != "https://" {
			person.WebsiteURL = &websiteURL
		}
		if givenName != "" {
			person.GivenName = &givenName
		}
		if err == nil {
			person.MemberFaves = &memberFaves
		}
		if len(altNames) != 0 {
			person.AlternateNames = altNames
		}
		if errT == nil {
			person.Birthday = &bdayT
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var familyName string
		name, ok := doc.Find("meta[property=og\\:title]").Attr("content")
		about := sel.Find("div.people-informantion-more.js-people-informantion-more").Text()
		darkText.Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Family name:" {
				familyName = strings.Trim(s.Nodes[0].NextSibling.Data, " ")
			}
		})
		if ok || name != "" {
			person.Name = &name
		}
		if about != "" {
			person.About = &about
		}
		if familyName != "" {
			person.FamilyName = &familyName
		}
	}()

	header := sel.Find("div.normal_header").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Nodes[0].LastChild.Data == "Voice Acting Roles"
	})
	tableRows := header.SiblingsFiltered("table").FilterFunction(func(i int, s *goquery.Selection) bool {
		return i == 0
	}).Children().Children()

	wg.Add(len(tableRows.Nodes))
	tableRows.Each(func(i int, s *goquery.Selection) {
		var role VARole
		person.Roles = append(person.Roles, &role)
		go func() {
			defer wg.Done()
			ScrapeVARole(&role, s)
		}()
	})

	wg.Wait()
	return person, nil
}

func ScrapeVARole(role *VARole, s *goquery.Selection) {
	role.Anime = &VAAnime{}
	role.Char = &VAChar{}

	animeLinkContainer := s.Find("div.picSurround")
	animeURL, ok := animeLinkContainer.Find("a").Attr("href")
	if ok || animeURL != "" {
		malID, err := strconv.Atoi(strings.Split(animeURL, "/")[4])
		if err == nil {
			role.Anime.MalID = &malID
		}
		role.Anime.URL = &animeURL
	}
	animeImgURL, _ := animeLinkContainer.Find("img").Attr("data-src")
	animeImgURL = strings.ReplaceAll(animeImgURL, "r/84x124/", "")
	animeImgURL = strings.ReplaceAll(animeImgURL, ".webp", ".jpg")
	role.Anime.ImageURL = &strings.Split(animeImgURL, "?")[0]

	character := s.Find("td.borderClass").Last().Find("div.picSurround")
	charURL, ok := character.Find("a").Attr("href")
	if ok || charURL != "" {
		malID, err := strconv.Atoi(strings.Split(charURL, "/")[4])
		if err == nil {
			role.Char.MalID = &malID
		}
		role.Char.URL = &charURL
	}
	charImg := character.Find("img")
	charName, ok := charImg.Attr("alt")
	if ok || charName != "" {
		role.Char.Name = &charName
	}

	imgURL, _ := charImg.Attr("data-src")
	imgURL = strings.ReplaceAll(imgURL, "r/84x124/", "")
	imgURL = strings.ReplaceAll(imgURL, ".webp", ".jpg")
	role.Char.ImageURL = &strings.Split(imgURL, "?")[0]

	role.Role = strings.Trim(s.Find("td").Nodes[2].LastChild.FirstChild.Data, " ")
}
