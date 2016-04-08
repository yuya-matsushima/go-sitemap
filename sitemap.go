package sitemap

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Index is a structure of <sitemapindex>
type Index struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemap []parts  `xml:"sitemap"`
}

// parts is a structure of <sitemap> in <sitemapindex>
type parts struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// Sitemap is a structure of <sitemap>
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URL     []URL    `xml:"url"`
}

// URL is a structure of <url> in <sitemap>
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float32 `xml:"priority"`
}

// fetch is page acquisition function
var fetch = func(URL string, timeout int64) ([]byte, error) {
	var body []byte

	res, err := http.Get(URL)
	if err != nil {
		return body, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	return body, err
}

// timeout is setting for fetch
var timeout int64 = 10

// Time interval to be used in Index.get
var interval = time.Second

// Get sitemap data from URL
func Get(url string) (Sitemap, error) {
	data, err := fetch(url, timeout)
	if err != nil {
		return Sitemap{}, err
	}

	index, indexErr := ParseIndex(data)
	sitemap, sitemapErr := Parse(data)

	if indexErr != nil && sitemapErr != nil {
		err = errors.New("URL is not a sitemap or sitemapindex")
		return Sitemap{}, err
	}

	if indexErr == nil {
		sitemap, err = index.get(data)
		if err != nil {
			return Sitemap{}, err
		}
	}

	return sitemap, err
}

// Get Sitemap data from sitemapindex file
func (s *Index) get(data []byte) (Sitemap, error) {
	index, err := ParseIndex(data)
	if err != nil {
		return Sitemap{}, err
	}

	var sitemap Sitemap
	for _, s := range index.Sitemap {
		time.Sleep(interval)
		data, err := fetch(s.Loc, timeout)
		if err != nil {
			return sitemap, err
		}

		err = xml.Unmarshal(data, &sitemap)
		if err != nil {
			return sitemap, err
		}
	}

	return sitemap, err
}

// Parse create Sitemap data from text
func Parse(data []byte) (Sitemap, error) {
	var sitemap Sitemap
	err := xml.Unmarshal(data, &sitemap)

	return sitemap, err
}

// ParseIndex create Index data from text
func ParseIndex(data []byte) (Index, error) {
	var index Index
	err := xml.Unmarshal(data, &index)

	return index, err
}

// SetInterval change Time interval to be used in Index.get
func SetInterval(time time.Duration) {
	interval = time
}

// SetFetch change fetch closure
func SetFetch(f func(url string, timeout int64) ([]byte, error)) {
	fetch = f
}
