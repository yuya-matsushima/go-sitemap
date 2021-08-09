package sitemap

import (
	"io/ioutil"
	"testing"
	"time"
)

// getTest is structure for test
type getTest struct {
	smapName string
	count    int
	hasErr   bool
	ErrStr   string
}

var getTests = []getTest{
	// sitemap.xml test
	{"sitemap.xml", 13, false, ""},
	// sitemap.xml is empty.
	{"empty_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: EOF"},
	// sitemap.xml is not exist.
	{"not_exist_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: EOF"},
	// sitemapindex.xml test
	{"sitemapindex.xml", 39, false, ""},
	// sitemapindex.xml is empty.
	{"empty_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: EOF"},
	// sitemapindex.xml is not exist.
	{"not_exist_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: EOF"},
	// sitemapindex.xml contains empty sitemap.xml
	{"contains_empty_sitemap_sitemapindex.xml", 0, true, "EOF"}, // TODO: fix error message
	// sitemapindex.xml contains sitemap.xml that is not exist.
	{"contains_not_exist_sitemap_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: EOF"},
}

func TestGet(t *testing.T) {
	server := testServer()
	defer server.Close()

	SetInterval(time.Nanosecond)

	for i, test := range getTests {
		data, err := Get(server.URL+"/"+test.smapName, nil)

		if test.hasErr {
			if err == nil {
				t.Errorf("%d: Get() should has error. expected:%s", i, test.ErrStr)
			}

			if err.Error() != test.ErrStr {
				t.Errorf("%d: Get() shoud return error. result:%s expected:%s", i, err.Error(), test.ErrStr)
			}
		} else {
			if err != nil {
				t.Errorf("%d: Get() should not has error. result: %s", i, err.Error())
			}
		}

		if test.count != len(data.URL) {
			t.Errorf("%d: Get() should return Sitemap.Url:%d expected: %d", i, len(data.URL), test.count)
		}
	}
}

func TestParse(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/sitemap.xml")
	smap, _ := Parse(data)

	if len(smap.URL) != 13 {
		t.Error("Parse() should return Sitemap.URL(13 length)")
	}
}

func TestParseIndex(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")
	idx, _ := ParseIndex(data)

	if len(idx.Sitemap) != 3 {
		t.Error("ParseIndex() should return Index.Sitemap(3 length)")
	}
}

func TestSetInterval(t *testing.T) {
	newInterval := 3 * time.Second
	SetInterval(newInterval)

	if interval != newInterval {
		t.Error("interval should be time.Minute")
	}

	if interval == time.Second {
		t.Error("interval should not be Default(time.Second)")
	}
}

func TestSetFetch(t *testing.T) {
	f := func(URL string, options interface{}) ([]byte, error) {
		var err error
		return []byte(URL), err
	}

	SetFetch(f)

	URL := "http://example.com"
	data, _ := fetch(URL, nil)

	if string(data) != URL {
		t.Error("fetch() should return " + URL)
	}
}

// func BenchmarkGetSitemap(b *testing.B) {
// 	server := testServer()
// 	defer server.Close()
//
// 	for i := 0; i < b.N; i++ {
// 		_, err := Get(server.URL+"/sitemap.xml", nil)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
//
// func BenchmarkGetSitemapIndex(b *testing.B) {
// 	server := testServer()
// 	defer server.Close()
//
// 	for i := 0; i < b.N; i++ {
// 		_, err := Get(server.URL+"/sitemapindex.xml", nil)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
//
// func BenchmarkParseSitemap(b *testing.B) {
// 	data, _ := ioutil.ReadFile("./testdata/sitemap.xml")
//
// 	for i := 0; i < b.N; i++ {
// 		_, err := Parse(data)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
//
// func BenchmarkParseSitemapIndex(b *testing.B) {
// 	data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")
//
// 	for i := 0; i < b.N; i++ {
// 		_, err := ParseIndex(data)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
