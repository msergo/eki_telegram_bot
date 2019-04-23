package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL             = "http://www.eki.ee/dict/evs/index.cgi?Q="
	cartSelector        = ".tervikart"
	paragraphSelector   = ".leitud_id"
	translationSelector = ".x_x"
)

// GetTranslations fnbad a
func GetTranslations() {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, "saama"))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(cartSelector).Each(func(i int, cart *goquery.Selection) {
		fmt.Println(cart.Find(paragraphSelector).Text())
		cart.Find(translationSelector).Each(func(i int, translation *goquery.Selection) {
			fmt.Println(translation.Text())
		})
	})
}

func main() {
	GetTranslations()
}
