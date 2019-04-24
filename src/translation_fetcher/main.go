package translation_fetcher

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL             = "http://www.eki.ee/dict/evs/index.cgi?Q="
	cartSelector        = ".tervikart"
	articleUseCaseSelector     = ".leitud_id"
	translationSelector = ".x_x"
)

// Meaning represents a sub-article
type Meaning struct {
	Word        string
	Translation string
	Examples []string
}

// Article represents each variant of word
type Article struct {
	Meanings      []Meaning
	ArticleHeader string
}

// GetTranslations fetches related article
func GetTranslations(word string) []Article {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, word))
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
	var articles []Article;
	// Find the review items
	doc.Find(cartSelector).Each(func(i int, article *goquery.Selection) {
		articleItem:= Article{}
		articleItem.ArticleHeader = article.Find(articleUseCaseSelector).Text()
		articleItem.Meanings = []Meaning{}
		article.Find(translationSelector).Each(func(i int, meaning *goquery.Selection) {
			meaningItem := Meaning{}
			meaningItem.Translation = meaning.Text()
			articleItem.Meanings = append(articleItem.Meanings, meaningItem)
		})
		articles = append(articles, articleItem)
	})

	return articles
}
