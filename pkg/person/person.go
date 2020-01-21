package person

import (
	"fmt"
	"github.com/deletescape/suzuha/internal/utils"
	scrp "github.com/deletescape/suzuha/pkg/scraper"
	"net/url"
	"sync"
)

type Person struct {
	URL *string `json:"url"`
}

var scraper scrp.Scraper

func ScrapePerson(id int) (*Person, error) {
	person := &Person{}
	requested := fmt.Sprintf("https://myanimelist.net/people/%d", id)
	sel, err := scraper.GetSelection(requested, "div#contentWrapper")
	if err != nil {
		return nil, err
	}
	_, err = url.Parse(requested)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	// Get person MAL URL
	wg.Add(1)
	go func() {
		defer wg.Done()
		navButtons := sel.Find(".horiznav_active")
		URL := utils.GetInfoLinkString(navButtons, "Details")
		if URL != "" {
			person.URL = &URL
		}
	}()

	// Get person image URL
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	wg.Wait()
	return person, nil
}
