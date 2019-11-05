package main

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	baseURL                = "http://www.eki.ee/dict/evs/index.cgi?Q="
	cartSelector           = ".tervikart"
	articleUseCaseSelector = ".leitud_id"
	translationSelector    = ".x_x[lang=\"ru\"]"
	exampleEstSelector     = ".x_n[lang=\"et\"]"
	exampleRusSelector     = ".x_qn[lang=\"ru\"]"
	grammarFormSelector    = ".x_mv[lang=\"et\"]"
)

// Meaning represents a sub-article
type Meaning struct {
	Word        string
	Translation string
	Examples    []string
}

// Article represents each variant of word
type Article struct {
	Meanings      []Meaning
	ArticleHeader string
}

// GetSingleArticle get preformatted translation
func GetSingleArticle(node *html.Node) (string, bool) {
	doc := goquery.NewDocumentFromNode(node)
	useCase := doc.Find(articleUseCaseSelector).Text()
	grammarForms := doc.Find(grammarFormSelector).Text()
	var translations []string
	// var examplesEst []string // temporary disable
	// var examples []string
	doc.Find(translationSelector).Each(func(i int, translation *goquery.Selection) {
		translations = append(translations, translation.Text())
	})

	// temporary disable
	// doc.Find(".x_n[lang=\"et\"], .x_qn[lang=\"ru\"], br").Each(func(i int, example *goquery.Selection) {
	// 	if example.Is(exampleEstSelector) {
	// 		examplesEst = append(examplesEst, example.Text())
	// 	} else if example.Is(exampleRusSelector) {
	// 		ex := examplesEst[0] + " - " + example.Text() // more than one example in russian is possible
	// 		examples = append(examples, ex)
	// 	} else if example.Is("br") {
	// 		examplesEst = examplesEst[:0] //clean buff
	// 	}
	// })
	if grammarForms == "" {
		return fmt.Sprintf("<b>%s</b>\r\n%s",
			useCase,
			strings.Join(translations, "\r\n"),
		), false
	}
	return fmt.Sprintf("<b>%s</b><i> (%s) </i>\r\n%s",
			useCase,
			grammarForms,
			strings.Join(translations, "\r\n"),
		),
		true

}

// fetchArticles fetches HTML page and returns as a collection of nodes
func fetchArticles(word string) *goquery.Document {
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
	return doc
}

// GetArticles fetches HTML page and extract separate word-related articles
func GetArticles(word string) []string {
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
	var articles []string

	doc.Find(cartSelector).Each(func(i int, page *goquery.Selection) {
		for i := 0; i < len(page.Nodes); i++ {
			article, isMainArticle := GetSingleArticle(page.Nodes[i])
			if isMainArticle {
				articles = append([]string{article}, articles...) //put main article to the first position
				continue
			}
			articles = append(articles, article)
		}
	})

	return articles
}
