package edgr

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/hodlgap/edgr/request"

	"golang.org/x/net/html/charset"

	"github.com/hodlgap/edgr/model"
)

var (
	secCompanyURL = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=%s&start=0&count=1&output=atom"
	dirRegex      = regexp.MustCompile(`<td><a href="(.*?)"><img`)
	urlRegex      = regexp.MustCompile(`.*<a href="(.*?)index.html"><img`)
)

// Filer models
// -----------------

// Company is a simple struct for a single company.
type Company struct {
	Name   string
	Symbol string
}

// rssFeed is the feed obj.
type rssFeed struct {
	Info secFilerInfo `xml:"company-info"`
}

type secFilerInfo struct {
	CIK     string `xml:"cik"`
	SIC     string `xml:"assigned-sic,omitempty"`
	SICDesc string `xml:"assigned-sic-desc,omitempty"`
	Name    string `xml:"conformed-name"`
}

// GetFiler gets a single filer from the SEC website based on symbol.
func GetFiler(symbol string) (*model.Filer, error) {
	// get the cik for each symbol.
	// tedious process...
	url := fmt.Sprintf(secCompanyURL, symbol)

	respBody, err := request.GetPage(url)
	if err != nil {
		return nil, err
	}

	var feed rssFeed
	decoder := xml.NewDecoder(strings.NewReader(respBody))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&feed); err != nil {
		return nil, errors.WithStack(err)
	}

	if feed.Info.CIK == "" {
		return nil, errors.New("no cik found in response data")
	}
	if feed.Info.Name == "" {
		return nil, errors.New("no name found in response data")
	}

	return &model.Filer{
		CIK:            feed.Info.CIK,
		Symbol:         symbol,
		SIC:            feed.Info.SIC,
		SICDescription: feed.Info.SICDesc,
		Name:           feed.Info.Name,
	}, nil
}

// Filings models
// -----------------

// SECFiling contains a single instance of an sec filing.
type SECFiling struct {
	Filing *model.Filing
	Docs   []*model.Document
}

// Filings methods
// -----------------

// GetFilings gets a list of filings for a single CIK.
func GetFilings(cik, formtype, stoptime string) (filings []SECFiling, err error) {
	var stop *time.Time
	if stoptime != "" {
		t, err := time.Parse("2006-01-02", stoptime)
		if err != nil {
			return filings, err
		}
		stop = &t
	}

	dirPage, err := request.GetPage("https://www.sec.gov/Archives/edgar/data/" + cik)
	if err != nil {
		return
	}

	urls := findListURLs(dirPage)

	for _, u := range urls {
		docsPage, getErr := request.GetPage(u)
		if getErr != nil {
			log.Errorf("couldnt find page: %s", getErr)
			continue
		}

		idxURL := findIdxURL(docsPage)
		if idxURL == "" {
			log.Error("couldnt regex idx url")
			continue
		}

		filing, buildErr := buildFiling(cik, idxURL)
		if buildErr != nil {
			log.Error(buildErr)
			continue
		}
		if formtype != "" {
			// check form type.
			if filing.Filing.FormType != formtype {
				continue
			}
		}

		if stop != nil {
			// check cutoff time.
			if filing.Filing.EdgarTime.Before(*stop) {
				return
			}
		}
		// Do stuff with the filing...
		filing.Filing.AllSymbols = []string{filing.Filing.Symbol}
		filings = append(filings, filing)
	}

	return
}

// findIdxURL parses the text document url out of the index page.
func findIdxURL(html string) string {
	matches := urlRegex.FindStringSubmatch(html)
	if matches == nil || len(matches) == 1 {
		log.Warn("could not find matches")
		return ""
	}
	return "https://sec.gov" + matches[1] + "index.html"
}

// findListURLs parses the list of idx urls out of the directory page.
func findListURLs(html string) []string {
	matches := dirRegex.FindAllStringSubmatch(html, -1)
	if matches == nil || len(matches) == 1 {
		log.Warn("could not find matches")
		return nil
	}

	urls := make([]string, len(matches))
	for i, m := range matches {
		urls[i] = "https://sec.gov" + m[1]
	}

	return urls
}
