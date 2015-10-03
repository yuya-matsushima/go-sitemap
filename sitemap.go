package sitemap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Index is a structure of <sitemapindex>
type Index struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemap []Parts  `xml:"sitemap"`
}

// Parts is a structure of <sitemap> in <sitemapindex>
type Parts struct {
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
