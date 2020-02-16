package main

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"regexp"
)

const (
	baseURL                = "http://www.eki.ee/dict/evs/index.cgi?Q="
	cartSelector           = ".tervikart"
	//articleUseCaseSelector = ".leitud_id" //TODO: update tests
	articleUseCaseSelector = ".m.x_m.m"
	translationSelector    = ".x_x[lang=\"ru\"]"
	exampleEstSelector     = ".x_n[lang=\"et\"]"
	exampleRusSelector     = ".x_qn[lang=\"ru\"]"
	grammarFormSelector    = ".mv.x_mv.mv[lang=\"et\"]"
)
var cleanupRegex, _ = regexp.Compile("[^a-zA-Z0-9]+")

func IsMatchingArticle(searchWord string, givenWord string) bool {
	a := strings.Split(givenWord, " ")
	var isMatch bool
	for i := range a {
		form := cleanupRegex.ReplaceAllString(a[i], "")
		if form == searchWord {
			isMatch = true
			break
		}
	}
	return isMatch
}

// GetSingleArticle get preformatted translation
func GetSingleArticle(searchWord string, node *html.Node) (string, bool) {
	doc := goquery.NewDocumentFromNode(node)
	useCase := doc.Find(articleUseCaseSelector).Text()
	grammarForms := doc.Find(grammarFormSelector).Text()
	//filter garbage
	if !IsMatchingArticle(searchWord, useCase) && !IsMatchingArticle(searchWord, grammarForms) {
		return "", false
	}
	var translations []string
	doc.Find(translationSelector).Each(func(i int, translation *goquery.Selection) {
		translations = append(translations, translation.Text())
	})

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


// GetArticles fetches HTML page and extract separate word-related articles
func GetArticles(searchWord string) []string {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, searchWord))
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
			article, isMainArticle := GetSingleArticle(searchWord, page.Nodes[i])
			if article == "" {
				continue
			}
			if isMainArticle {
				articles = append([]string{article}, articles...) //put main article to the first position
				continue
			}
			articles = append(articles, article)
		}
	})

	return articles
}
