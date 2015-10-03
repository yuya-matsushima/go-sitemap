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

// Parts is a structure of <sitemap> in <sitemapindex>
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

var fetch = func(url string) ([]byte, error) {
	var body []byte

	res, err := http.Get(url)
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

// Get sitemap data from URL
func GetSitemap(url string) (Sitemap, error) {
	var index Index
	var sitemap Sitemap

	data, err := fetch(url)
	if err != nil {
		return sitemap, err
	}

	indexErr := xml.Unmarshal(data, &index)
	sitemapErr := xml.Unmarshal(data, &sitemap)

	if indexErr != nil && sitemapErr != nil {
		err = errors.New("URL is not a sitemap or sitemapindex")
		return Sitemap{}, err
	}

	if indexErr == nil {
		sitemap, err = index.get(data)
		if err != nil {
			return sitemap, err
		}
	}

	return sitemap, err
}

// Get Sitemap data from sitemapindex file
func (s *Index) get(data []byte) (Sitemap, error) {
	var index Index
	var sitemap Sitemap

	err := xml.Unmarshal(data, &index)
	if err != nil {
		return Sitemap{}, err
	}

	for _, s := range index.Sitemap {
		time.Sleep(time.Second) // TODO: sleep time will be option.
		data, err := fetch(s.Loc)
		if err != nil {
			return sitemap, err
		}
		xml.Unmarshal(data, &sitemap)
	}

	return sitemap, err
}
