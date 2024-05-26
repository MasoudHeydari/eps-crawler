package google

import (
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/gocolly/colly"
	"github.com/karust/openserp/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Google struct {
	core.Browser
	core.SearchEngineOptions
	collector                                  *colly.Collector
	findNumRgxp, findPhoneRgxp, findEmailRegex *regexp.Regexp
	client                                     *http.Client
}

func New(browser core.Browser, opts core.SearchEngineOptions) *Google {
	gogl := Google{
		Browser:   browser,
		collector: colly.NewCollector(),
		client:    &http.Client{Timeout: 5 * time.Second},
	}
	opts.Init()
	gogl.SearchEngineOptions = opts

	gogl.findNumRgxp = regexp.MustCompile("\\d")
	gogl.findPhoneRgxp = regexp.MustCompile(`\b\(?\+?\d{1,2}[\s\-]?\)?\(?\d{3}\)?[\s\-\.]?\d{3}[\s\-\.]?\d{4}\b`) // regexp.MustCompile(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$`) //`\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`)
	gogl.findEmailRegex = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	return &gogl
}

func (gogl *Google) Name() string {
	return "google"
}

func (gogl *Google) GetRateLimiter() *rate.Limiter {
	ratelimit := rate.Every(gogl.GetRatelimit())
	return rate.NewLimiter(ratelimit, gogl.RateBurst)
}

func (gogl *Google) findTotalResults(page *rod.Page) (int, error) {
	resultsStats, err := page.Timeout(gogl.GetSelectorTimeout()).Search("div#result-stats")
	if err != nil {
		return 0, errors.New("Result stats not found: " + err.Error())
	}

	stats, err := resultsStats.First.Text()
	if err != nil {
		return 0, errors.New("Cannot extract result stats text: " + err.Error())
	}

	// Escape moment with `seconds` and extract digits
	allNums := gogl.findNumRgxp.FindAllString(stats[:len(stats)-15], -1)
	stats = strings.Join(allNums, "")

	total, err := strconv.Atoi(stats)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (gogl *Google) isCaptcha(page *rod.Page) bool {
	_, err := page.Timeout(gogl.GetSelectorTimeout()).Search("form#captcha-form")
	if err != nil {
		return false
	}
	return true
}

func (gogl *Google) preparePage(page *rod.Page) {
	// Remove "similar queries" lists
	page.Eval(";(() => { document.querySelectorAll(`div[data-initq]`).forEach( el => el.remove());  })();")
}

func (gogl *Google) Search(query core.Query) ([]core.SearchResult, error) {
	logrus.Tracef("Start Google search, query: %+v", query)

	var searchResults []core.SearchResult

	// Build URL from query struct to open in browser
	url, err := BuildURL(query)
	if err != nil {
		return nil, err
	}

	page := gogl.Navigate(url)
	gogl.preparePage(page)

	results, err := page.Timeout(gogl.Timeout).Search("div[data-hveid][data-ved][lang], div[data-surl][jsaction]")
	if err != nil {
		defer page.Close()
		logrus.Errorf("Cannot parse search results: %s", err)
		return nil, core.ErrSearchTimeout
	}

	// Check why no results, maybe captcha?
	if results == nil {
		defer page.Close()

		if gogl.isCaptcha(page) {
			logrus.Errorf("Google captcha occurred during: %s", url)
			return nil, core.ErrCaptcha
		}
		return nil, err
	}

	totalResults, err := gogl.findTotalResults(page)
	if err != nil {
		logrus.Errorf("Error capturing total results: %v", err)
	}
	logrus.Infof("%d total results found", totalResults)

	resultElements, err := results.All()
	if err != nil {
		return nil, err
	}

	for i, r := range resultElements {
		// Get URL
		link, err := r.Element("a")
		if err != nil {
			continue
		}
		linkText, err := link.Property("href")
		if err != nil {
			logrus.Error("No `href` tag found")
		}

		// Get title
		titleTag, err := link.Element("h3")
		if err != nil {
			logrus.Error("No `h3` tag found")
			continue
		}

		title, err := titleTag.Text()
		if err != nil {
			logrus.Error("Cannot extract text from title")
			title = "No title"
		}

		// Get description
		// doesn't catch all
		descTag, err := r.Element(`div[data-sncf~="1"]`)
		desc := ""
		if err != nil {
			logrus.Trace(`No description 'div[data-sncf~="1"]' tag found`)
		} else {
			desc = descTag.MustText()
		}

		// extract contact-info
		var contactInfo *core.ContactInfo
		contactInfo, err = gogl.extractContactInfo(linkText.String())
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}

		// extract key-words
		var keyWords []string
		keyWords, err = gogl.extractKeywords(linkText.String())
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}

		gR := core.SearchResult{
			Rank:        i + 1,
			URL:         linkText.String(),
			Title:       title,
			ContactInfo: contactInfo,
			KeyWords:    keyWords,
			Description: desc,
		}
		searchResults = append(searchResults, gR)
	}

	if !gogl.Browser.LeavePageOpen {
		err = page.Close()
		if err != nil {
			logrus.Error(err)
		}
	}

	return searchResults, nil
}

func (gogl *Google) extractKeywords(path string) ([]string, error) {
	r := make(map[string]struct{})
	gogl.collector.OnHTML("h1, h2", func(e *colly.HTMLElement) {
		r[e.Text] = struct{}{}
	})
	err := gogl.collector.Visit(path)
	if err != nil {
		return nil, fmt.Errorf("extractKeywords: %w", err)
	}
	keyWords := make([]string, 0, 5) // TODO: use proper capacity
	for k := range r {
		keyWords = append(keyWords, k)
	}
	return keyWords, nil
}

func (gogl *Google) extractContactInfo(path string) (*core.ContactInfo, error) {
	phonesMap := make(map[string]struct{})
	emailsMap := make(map[string]struct{})
	resp, err := gogl.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("extractContactInfo.Get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("extractContactInfo.ReadAll: %w", err)
	}
	phonesArray := gogl.findPhoneRgxp.FindAllString(string(body), -1)
	emailsArray := gogl.findEmailRegex.FindAllString(string(body), -1)
	for _, phone := range phonesArray {
		phonesMap[phone] = struct{}{}
	}

	for _, email := range emailsArray {
		emailsMap[email] = struct{}{}
	}
	phones := make([]string, 0, len(phonesMap))
	emails := make([]string, 0, len(emailsMap))

	for k := range emailsMap {
		emails = append(emails, k)
	}

	for k := range phonesMap {
		phones = append(phones, k)
	}

	return &core.ContactInfo{Phones: phones, Emails: emails}, nil
}
