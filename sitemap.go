package sitemap

import (
	"encoding/xml"
	"fmt"
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

var (
	// fetch is page acquisition function
	fetch = func(URL string, options interface{}) ([]byte, error) {
		var body []byte

		res, err := http.Get(URL)
		if err != nil {
			return body, err
		}
		defer res.Body.Close()

		return ioutil.ReadAll(res.Body)
	}

	// Time interval to be used in Index.get
	interval = time.Second
)

// Get sitemap data from URL
func Get(URL string, options interface{}) (Sitemap, error) {
	data, err := fetch(URL, options)
	if err != nil {
		return Sitemap{}, err
	}

	idx, idxErr := ParseIndex(data)
	smap, smapErr := Parse(data)

	if idxErr != nil && smapErr != nil {
		if idxErr != nil {
			err = idxErr
		} else {
			err = smapErr
		}
		return Sitemap{}, fmt.Errorf("URL is not a sitemap or sitemapindex.: %v", err)
	} else if idxErr != nil {
		return smap, nil
	}

	smap, err = idx.get(options)
	if err != nil {
		return Sitemap{}, err
	}

	return smap, nil
}

// Get Sitemap data from sitemapindex file
func (idx *Index) get(options interface{}) (Sitemap, error) {
	var smap Sitemap

	for _, s := range idx.Sitemap {
		time.Sleep(interval)
		data, err := fetch(s.Loc, options)
		if err != nil {
			return smap, fmt.Errorf("failed to retrieve %s in sitemapindex.xml.: %v", s.Loc, err)
		}

		err = xml.Unmarshal(data, &smap)
		if err != nil {
			return smap, fmt.Errorf("failed to parse %s in sitemapindex.xml.: %v", s.Loc, err)
		}
	}

	return smap, nil
}

// Parse create Sitemap data from text
func Parse(data []byte) (Sitemap, error) {
	var smap Sitemap
	if len(data) == 0 {
		return smap, fmt.Errorf("sitemap.xml is empty.")
	}

	err := xml.Unmarshal(data, &smap)
	return smap, err
}

// ParseIndex create Index data from text
func ParseIndex(data []byte) (Index, error) {
	var idx Index
	if len(data) == 0 {
		return idx, fmt.Errorf("sitemapindex.xml is empty.")
	}

	err := xml.Unmarshal(data, &idx)
	return idx, err
}

// SetInterval change Time interval to be used in Index.get
func SetInterval(time time.Duration) {
	interval = time
}

// SetFetch change fetch closure
func SetFetch(f func(URL string, options interface{}) ([]byte, error)) {
	fetch = f
}
