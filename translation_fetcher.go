package main

import (
	"fmt"
	"net/http"

	"strings"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/getsentry/sentry-go"
	"golang.org/x/net/html"
)

const (
	baseURLEstRus    = "http://www.eki.ee/dict/evs/index.cgi?Q="
	baseURLRusEst = "http://www.eki.ee/dict/ves/index.cgi?Q="
	baseURLEstUkr = "http://www.eki.ee/dict/ukraina/index.cgi?Q="

	cartSelector = ".tervikart"
	//articleUseCaseSelector = ".leitud_id" //TODO: update tests
	articleUseCaseSelector    = ".m.x_m.m"
	articleUseCaseSelectorRus = ".ms.leitud_id"
	translationSelector       = ".x_x[lang=\"ru\"]"
	exampleEstSelector        = ".x_n[lang=\"et\"]"
	exampleRusSelector        = ".x_qn[lang=\"ru\"]"
	grammarFormSelector       = ".mv.x_mv.mv[lang=\"et\"]"
)

var cleanupRegex, _ = regexp.Compile("[^\\p{L}]+")

func IsMatchingArticle(searchWord string, givenWord string) bool {
	a := strings.Split(givenWord, " ")
	var isMatch = false
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
	var useCase string
	if isCyrillicScript(searchWord) { // TODO: refactor
		text := doc.Text()
		text = strings.Replace(text, ";", "\r\n", -1)
		useCase = doc.Find(articleUseCaseSelectorRus).Text()
		if !IsMatchingArticle(searchWord, useCase) {
			return "", false
		}

		return text, false
	}
	useCase = doc.Find(articleUseCaseSelector).Text()
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
	var url string

	if isCyrillicScript(searchWord) {
		url = baseURLRusEst 
	} else {
		url = baseURLEstUkr
	}
	res, err := http.Get(fmt.Sprintf("%s%s", url, searchWord))
	captureErrorIfNotNull(err)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		sentry.CaptureException(fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	captureErrorIfNotNull(err)
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

func isCyrillicScript(searchWord string) bool {
	var rxCyrillic = regexp.MustCompile("^[\u0400-\u04FF\u0500-\u052F]+$")
	return rxCyrillic.MatchString(searchWord)
}
