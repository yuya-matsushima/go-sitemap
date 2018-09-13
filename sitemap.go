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
var fetch = func(URL string, options interface{}) ([]byte, error) {
	var body []byte

	res, err := http.Get(URL)
	if err != nil {
		return body, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// Time interval to be used in Index.get
var interval = time.Second

// Get sitemap data from URL
func Get(URL string, options interface{}) (Sitemap, error) {
	data, err := fetch(URL, options)
	if err != nil {
		return Sitemap{}, err
	}

	idx, idxErr := ParseIndex(data)
	smap, smapErr := Parse(data)

	if idxErr != nil && smapErr != nil {
		return Sitemap{}, errors.New("URL is not a sitemap or sitemapindex")
	} else if idxErr != nil {
		return smap, nil
	}

	smap, err = idx.get(data, options)
	if err != nil {
		return Sitemap{}, err
	}

	return smap, nil
}

// Get Sitemap data from sitemapindex file
func (s *Index) get(data []byte, options interface{}) (Sitemap, error) {
	idx, err := ParseIndex(data)
	if err != nil {
		return Sitemap{}, err
	}

	var smap Sitemap
	for _, s := range idx.Sitemap {
		time.Sleep(interval)
		data, err := fetch(s.Loc, options)
		if err != nil {
			return smap, err
		}

		err = xml.Unmarshal(data, &smap)
		if err != nil {
			return smap, err
		}
	}

	return smap, err
}

// Parse create Sitemap data from text
func Parse(data []byte) (smap Sitemap, err error) {
	err = xml.Unmarshal(data, &smap)
	return
}

// ParseIndex create Index data from text
func ParseIndex(data []byte) (idx Index, err error) {
	err = xml.Unmarshal(data, &idx)
	return
}

// SetInterval change Time interval to be used in Index.get
func SetInterval(time time.Duration) {
	interval = time
}

// SetFetch change fetch closure
func SetFetch(f func(URL string, options interface{}) ([]byte, error)) {
	fetch = f
}
