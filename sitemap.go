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

// Asp is a structure of <AspNet Sitemap-File>
type Asp struct {
	XMLName      xml.Name `xml:"siteMap"`
	SitemapNodes []Node   `xml:"siteMapNode"`
}

// Node is a structure of <sitemapNode> in <AspNet Sitemap-File>
type Node struct {
	URL          string `xml:"url,attr"`
	SitemapNodes []Node `xml:"siteMapNode"`
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
	asp, aspErr := ParseAsp(data)
	smap, smapErr := Parse(data)

	if idxErr != nil && smapErr != nil && aspErr != nil {
		return Sitemap{}, errors.New("URL is not a sitemap or sitemapindex or asp sitemap")
	} else if idxErr != nil {
		if aspErr != nil {
			return smap, nil
		}
		return asp.get(), nil
	}

	return idx.get(data, options)
}

// Get Sitemap data from sitemapindex file
func (idx *Index) get(data []byte, options interface{}) (Sitemap, error) {
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

	return smap, nil
}

// Get Sitemap data from asp sitemap file
func (a *Asp) get() Sitemap {
	var smap Sitemap
	addToSitemap(&smap, a.SitemapNodes)

	return smap
}

func addToSitemap(smap *Sitemap, nodes []Node) {
	for _, s := range nodes {
		smap.URL = append(smap.URL, URL{Loc: s.URL})
		if len(s.SitemapNodes) != 0 {
			addToSitemap(smap, s.SitemapNodes)
		}
	}
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

// ParseAsp create Asp data from text
func ParseAsp(data []byte) (asp Asp, err error) {
	err = xml.Unmarshal(data, &asp)
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
