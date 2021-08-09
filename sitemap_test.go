package sitemap

import (
	"io/ioutil"
	"testing"
	"time"
)

// getTest is structure for test
type getTest struct {
	smapName string
	isNil    bool
	count    int
}

var getTests = []getTest{
	// normal test
	{"sitemap.xml", true, 13},
	// This sitemap.xml is not exist.
	{"empty.xml", false, 0},
	// sitemap index test
	{"sitemapindex.xml", true, 39},
}

func TestGet(t *testing.T) {
	server := testServer()
	defer server.Close()

	SetInterval(time.Nanosecond)

	for i, test := range getTests {
		data, err := Get(server.URL+"/"+test.smapName, nil)

		if test.isNil == true && err != nil {
			t.Errorf("test:%d Get() should not has error:%s", i, err.Error())
		} else if test.isNil == false && err == nil {
			t.Errorf("test:%d Get() should has error", i)
		}

		if test.count != len(data.URL) {
			t.Errorf("test:%d Get() should return Sitemap.Url:%d actual: %d", i, test.count, len(data.URL))
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

func BenchmarkGetSitemap(b *testing.B) {
	server := testServer()
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, err := Get(server.URL+"/sitemap.xml", nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetSitemapIndex(b *testing.B) {
	server := testServer()
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, err := Get(server.URL+"/sitemapindex.xml", nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseSitemap(b *testing.B) {
	data, _ := ioutil.ReadFile("./testdata/sitemap.xml")

	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseSitemapIndex(b *testing.B) {
	data, _ := ioutil.ReadFile("./testdata/sitemapindex.xml")

	for i := 0; i < b.N; i++ {
		_, err := ParseIndex(data)
		if err != nil {
			b.Error(err)
		}
	}
}
