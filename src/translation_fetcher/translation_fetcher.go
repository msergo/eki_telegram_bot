package translation_fetcher

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"strings"
	"golang.org/x/net/html"
)

const (
	baseURL                = "http://www.eki.ee/dict/evs/index.cgi?Q="
	cartSelector           = ".tervikart"
	articleUseCaseSelector = ".leitud_id"
	translationSelector    = ".x_x[lang=\"ru\"]"
	exampleEstSelector     = ".x_n[lang=\"et\"]"
	exampleRusSelector     = ".x_qn[lang=\"ru\"]"
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

func GetSingleArticle(node *html.Node) string {
	doc := goquery.NewDocumentFromNode(node)

	useCase := doc.Find(articleUseCaseSelector).Text()
	var translations []string
	var examplesEst []string
	var examples []string
	doc.Find(translationSelector).Each(func(i int, translation *goquery.Selection) {
		translations = append(translations, translation.Text())
	})
	doc.Find(".x_n[lang=\"et\"], .x_qn[lang=\"ru\"], br").Each(func(i int, example *goquery.Selection) {
		if example.Is(exampleEstSelector) {
			examplesEst = append(examplesEst, example.Text())
		} else if example.Is(exampleRusSelector) {
			ex := examplesEst[0] + " - " + example.Text() // more than one example in russian is possible
			examples = append(examples, ex)
		} else if example.Is("br") {
			examplesEst = examplesEst[:0] //clean buff
		}
	})
	return fmt.Sprintf("<b>%s</b>\r\n"+
		//"%s\r\n"+ // remove examples for now
		"%s",
		useCase,
		strings.Join(translations, "\r\n"),
		//strings.Join(examples, "\r\n"), // remove examples for now
	)
}
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
	var articles []string
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(cartSelector).Each(func(i int, page *goquery.Selection) {
		for i := 0; i < len(page.Nodes); i++ {
			articles = append(articles, GetSingleArticle(page.Nodes[i]))
		}
	})

	return articles
}