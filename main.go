package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Fish struct {
	Species string `json:"species"`
	Status  string `json:"status"`
	Year    string `json:"year"`
	Region  string `json:"region"`
}

func Webscrapper() {
	fishes := []Fish{}
	space := regexp.MustCompile(`\s+`)
	c := colly.NewCollector(
		colly.AllowedDomains("www.fisheries.noaa.gov"),
	)

	c.OnHTML("div.species-directory__species--8col", func(element *colly.HTMLElement) {
		species := element.DOM

		fish := Fish{
			Species: species.Find("div.species-directory__species-title--name").Text(),
			Status:  space.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-status").Text()), " "),
			Year:    space.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-year").Text()), " "),
			Region:  space.ReplaceAllString(strings.TrimSpace(species.Find("div.species-directory__species-status-row").Find("div.species-directory__species-region").Text()), " "),
		}
		fishes = append(fishes, fish)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.fisheries.noaa.gov/species-directory/threatened-endangered?title=&species_category=any&species_status=any&regions=all&items_per_page=all&sort=")

	writeJson(fishes)
}

func writeJson(fishes []Fish) {
	f, err := json.MarshalIndent(fishes, "", " ")
	if err != nil {
		log.Fatal(err)
		return
	}

	_ = ioutil.WriteFile("endangered.json", f, 0644)
}

func main() {
	Webscrapper()
}
