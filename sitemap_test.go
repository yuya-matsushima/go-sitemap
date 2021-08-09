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
	t.Run("sitemap.xml exists", func(t *testing.T) {
		data, _ := ioutil.ReadFile("./testdata/sitemap.xml")
		smap, err := Parse(data)

		if err != nil {
			t.Errorf("Parse() should not return error. result:%v", err)
		}

		if len(smap.URL) != 13 {
			t.Errorf("Parse() should return Sitemap.URL. result:%d expected:%d", 13, len(smap.URL))
		}
	})

	t.Run("sitemap.xml not exists", func(t *testing.T) {
		smap, err := Parse([]byte{})

		if err.Error() != "sitemap.xml is empty." {
			t.Errorf("Parse() should return error. result:%s expected:%s", err.Error(), "sitemap.xml is empty.")
		}

		if len(smap.URL) != 0 {
			t.Errorf("Parse() should return Sitemap.URL. result:%d expected:%d", 0, len(smap.URL))
		}
	})
}

func TestParseIndex(t *testing.T) {
	t.Run("sitemapindex.xml exists", func(t *testing.T) {
		data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")
		idx, err := ParseIndex(data)

		if err != nil {
			t.Errorf("ParseIndex() should not return error. result:%v", err)
		}

		if len(idx.Sitemap) != 3 {
			t.Errorf("ParseIndex() should return Sitemap. result:%d expected:%d", 3, len(idx.Sitemap))
		}
	})

	t.Run("sitemapinde.xml not exists", func(t *testing.T) {
		idx, err := ParseIndex([]byte{})

		if err.Error() != "sitemapindex.xml is empty." {
			t.Errorf("ParseIndex() should not return error. result:%s expected:%s", err.Error(), "sitemapindex.xml is empty.")
		}

		if len(idx.Sitemap) != 0 {
			t.Errorf("ParseIndex() should return Sitemap. result:%d expected:%d", 0, len(idx.Sitemap))
		}
	})
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
