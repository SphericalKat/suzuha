package person

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deletescape/suzuha/internal/utils"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Person struct {
	URL            *string  `json:"url"`
	ImageURL       *string  `json:"image_url"`
	WebsiteURL     *string  `json:"website_url"`
	GivenName      *string  `json:"given_name"`
	FamilyName     *string  `json:"family_name"`
	About          *string  `json:"about"`
	MalID          *int     `json:"mal_id"`
	MemberFaves    *int     `json:"member_favorites"`
	AlternateNames []string `json:"alternate_names"`
}

var scraper scrp.Scraper

func ScrapePerson(id int) (*Person, error) {
	person := &Person{MalID: &id}
	requested := fmt.Sprintf("https://myanimelist.net/people/%d", id)
	sel, err := scraper.GetSelection(requested, "div#contentWrapper")
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
		var familyName string
		s := utils.GetInfoLinks(darkText, "Website:")
		websiteURL, _ := s.Attr("href")
		givenName := utils.GetInfo(darkText, "Given name:")
		memberFaves, err := strconv.Atoi(strings.ReplaceAll(utils.GetInfo(darkText, "Member Favorites:"), ",", ""))
		about := sel.Find("div.people-informantion-more.js-people-informantion-more").Text()
		altNames := strings.Split(utils.GetInfo(darkText, "Alternate names:"), ", ")
		darkText.Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Family name:" {
				familyName = strings.Trim(s.Nodes[0].NextSibling.Data, " ")
			}
		})

		if websiteURL != "" && websiteURL != "http://" && websiteURL != "https://" {
			person.WebsiteURL = &websiteURL
		}
		if givenName != "" {
			person.GivenName = &givenName
		}
		if familyName != "" {
			person.FamilyName = &familyName
		}
		if err == nil {
			person.MemberFaves = &memberFaves
		}
		if about != "" {
			person.About = &about
		}
		if len(altNames) != 0 {
			person.AlternateNames = altNames
		}

	}()

	//
	wg.Add(1)
	go func() {
		defer wg.Done()

	}()

	wg.Wait()
	return person, nil
}
