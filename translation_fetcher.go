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
	baseURLEstRus = "http://www.eki.ee/dict/evs/index.cgi?Q="
	baseURLRusEst = "http://www.eki.ee/dict/ves/index.cgi?Q="
	baseURLEstUkr = "http://www.eki.ee/dict/ukraina/index.cgi?Q="

	cartSelector              = ".tervikart"
	articleUseCaseSelectorEst = ".m.x_m.m"
	articleUseCaseSelectorRus = ".ms.leitud_id"
	translationSelectorEstRus = ".x_x[lang=\"ru\"]"
	exampleEstSelector        = ".x_n[lang=\"et\"]"
	exampleRusSelector        = ".x_qn[lang=\"ru\"]"
	grammarFormSelectorEst    = ".mv.x_mv.mv[lang=\"et\"]"

	articleUseCaseSelectorEstUkr = ".k_m.m"
	translationSelectorEstUkr    = ".k_x.x[lang=\"uk\"]"
	grammarFormSelectorEstUkr    = ".k_mv.mv[lang=\"et\"]"
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

func extractTranslationsFromDoc(doc *goquery.Document, searchWord string, useCase string, translationSelector string, grammarFormSelector string) []string {
	var translations []string
	grammarForms := doc.Find(grammarFormSelector).Text()
	// filter garbage
	if !IsMatchingArticle(searchWord, useCase) && !IsMatchingArticle(searchWord, grammarForms) {
		return translations
	}
	doc.Find(translationSelector).Each(func(i int, translation *goquery.Selection) {
		translations = append(translations, translation.Text())
	})
	return translations
}

func GetSingleArticleWithDirection(searchWord string, node *html.Node, direction string) (string, bool) {
	var articeUseCaseSelector string
	var grammarFormSelector string
	var translationSelector string
	var useCase string

	doc := goquery.NewDocumentFromNode(node)

	// TODO: add constants for directions
	if direction == "est-ukr" {
		articeUseCaseSelector = articleUseCaseSelectorEstUkr
		grammarFormSelector = grammarFormSelectorEstUkr
		translationSelector = translationSelectorEstUkr
	} else if direction == "est-rus" {
		articeUseCaseSelector = articleUseCaseSelectorEst
		grammarFormSelector = grammarFormSelectorEst
		translationSelector = translationSelectorEstRus
	} else if direction == "rus-est" {
		articeUseCaseSelector = articleUseCaseSelectorRus
	}

	useCase = doc.Find(articeUseCaseSelector).Text()

	if !IsMatchingArticle(searchWord, useCase) {
		return "", false
	}

	if direction == "rus-est" {
		text := doc.Text()
		text = strings.Replace(text, ";", "\r\n", -1)

		return text, false
	}

	grammarForms := doc.Find(grammarFormSelector).Text()
	translations := extractTranslationsFromDoc(doc, searchWord, useCase, translationSelector, grammarFormSelector)

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

// ParseArticles extracts article from the HTML node doc
func ParseArticles(doc *goquery.Document, searchWord string, translationDirection string) []string {
	var articles []string

	// Check if the search word is found at all
	foundArticlesInfo := doc.Find("P.inf").Text()

	// TODO: define separate filter function
	// TODO: check cases when found articles are not what we are looking for
	if strings.Contains(foundArticlesInfo, "Küsitud kujul või valitud artikli osast otsitut ei leitud, kasutan laiendatud otsingut") {
		return articles
	}

	doc.Find(cartSelector).Each(func(i int, page *goquery.Selection) {
		for i := 0; i < len(page.Nodes); i++ {
			article, isMainArticle := GetSingleArticleWithDirection(searchWord, page.Nodes[i], translationDirection)
			if article == "" {
				continue
			}

			if translationDirection == "est-ukr" {
				// prepend ukrainian flag emoji
				article = "🇺🇦 " + article
			}

			if isMainArticle {
				articles = append([]string{article}, articles...) // Put main article to the first position
				continue
			}
			articles = append(articles, article)
		}
	})

	return articles
}

// GetArticles fetches HTML page and extract separate word-related articles
func GetArticles(searchWord string) []string {
	// Request the HTML page.
	var urls []string
	var translationDirections []string

	if isCyrillicScript(searchWord) {
		urls = []string{baseURLRusEst}
		translationDirections = []string{"rus-est"}
	} else {
		urls = []string{baseURLEstRus, baseURLEstUkr}
		translationDirections = []string{"est-rus", "est-ukr"}
	}

	var articles []string
	for i := 0; i < len(urls); i++ {
		url := urls[i]
		translationDirection := translationDirections[i]

		res, err := http.Get(fmt.Sprintf("%s%s", url, searchWord))
		captureErrorIfNotNull(err)
		defer res.Body.Close()

		if res.StatusCode != 200 {
			sentry.CaptureException(fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status))
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		captureErrorIfNotNull(err)

		articlesPerDirection := ParseArticles(doc, searchWord, translationDirection)
		articles = append(articles, articlesPerDirection...)
	}

	return articles
}

func isCyrillicScript(searchWord string) bool {
	var rxCyrillic = regexp.MustCompile("^[\u0400-\u04FF\u0500-\u052F]+$")
	return rxCyrillic.MatchString(searchWord)
}
