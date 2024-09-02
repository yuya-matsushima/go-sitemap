package sitemap

import (
	"os"
	"strings"
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
	{"empty_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/empty_sitemap.xml"},
	// sitemap.xml is not exist.
	{"not_exist_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/not_exist_sitemap.xml"},
	// sitemapindex.xml test
	{"sitemapindex.xml", 39, false, ""},
	// sitemapindex.xml is empty.
	{"empty_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/empty_sitemapindex.xml"},
	// sitemapindex.xml is not exist.
	{"not_exist_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/not_exist_sitemapindex.xml"},
	// sitemapindex.xml contains empty sitemap.xml
	{"contains_empty_sitemap_sitemapindex.xml", 0, true, "failed to parse http://HOST/empty_sitemap.xml in sitemapindex.xml: EOF"},
	// sitemapindex.xml contains sitemap.xml that is not exist.
	{"contains_not_exist_sitemap_sitemapindex.xml", 0, true, "failed to parse http://HOST/not_exist_sitemap.xml in sitemapindex.xml: EOF"},
}

func TestGet(t *testing.T) {
	server := testServer()
	defer server.Close()

	SetInterval(time.Nanosecond)

	for i, test := range getTests {
		data, err := Get(server.URL+"/"+test.smapName, nil)

		// replace HOST in Error Message
		errMsg := test.ErrStr
		if strings.Contains(errMsg, "HOST") {
			errMsg = strings.Replace(errMsg, "http://HOST", server.URL, 1)
		}

		if test.hasErr {
			if err == nil {
				t.Errorf("%d: Get() should has error. expected:%s", i, errMsg)
			}

			if err.Error() != errMsg {
				t.Errorf("%d: Get() shoud return error. result:%s expected:%s", i, err.Error(), errMsg)
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

var forceGetTests = []getTest{
	// sitemap.xml test
	{"sitemap.xml", 13, false, ""},
	// sitemap.xml is empty.
	{"empty_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/empty_sitemap.xml"},
	// sitemap.xml is not exist.
	{"not_exist_sitemap.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/not_exist_sitemap.xml"},
	// sitemapindex.xml test
	{"sitemapindex.xml", 39, false, ""},
	// sitemapindex.xml is empty.
	{"empty_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/empty_sitemapindex.xml"},
	// sitemapindex.xml is not exist.
	{"not_exist_sitemapindex.xml", 0, true, "URL is not a sitemap or sitemapindex: http://HOST/not_exist_sitemapindex.xml"},
	// sitemapindex.xml contains empty sitemap.xml
	{"contains_empty_sitemap_sitemapindex.xml", 13, false, ""},
	// sitemapindex.xml contains sitemap.xml that is not exist.
	{"contains_not_exist_sitemap_sitemapindex.xml", 13, false, ""},
}

func TestForceGet(t *testing.T) {
	server := testServer()
	defer server.Close()

	SetInterval(time.Nanosecond)

	for i, test := range forceGetTests {
		data, err := ForceGet(server.URL+"/"+test.smapName, nil)

		// replace HOST in Error Message
		errMsg := test.ErrStr
		if strings.Contains(errMsg, "HOST") {
			errMsg = strings.Replace(errMsg, "http://HOST", server.URL, 1)
		}

		if test.hasErr {
			if err == nil {
				t.Errorf("%d: Get() should has error. expected:%s", i, errMsg)
			}

			if err.Error() != errMsg {
				t.Errorf("%d: Get() shoud return error. result:%s expected:%s", i, err.Error(), errMsg)
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

func TestReadSitemap(t *testing.T) {
	t.Run("sitemap.xml exists", func(t *testing.T) {
		path := "./testdata/sitemap.xml"
		smap, err := ReadSitemap(path)

		if err != nil {
			t.Errorf("ReadSitemap() should not return error. result:%v", err)
		}

		if len(smap.URL) != 13 {
			t.Errorf("ReadSitemap() should return Sitemap.URL. result:%d expected:%d", 13, len(smap.URL))
		}
	})

	t.Run("sitemap.xml not exists", func(t *testing.T) {
		path := "./testdata/not_exist_sitemap.xml"
		smap, err := ReadSitemap(path)

		errText := "file not found ./testdata/not_exist_sitemap.xml"
		if err.Error() != errText {
			t.Errorf("ReadSitemap() should return error. result:%s expected:%s", err.Error(), errText)
		}

		if len(smap.URL) != 0 {
			t.Errorf("ReadSitemap() should not return Sitemap.URL. result:%d expected:%d", 0, len(smap.URL))
		}
	})
}

func TestReadSitemapIndex(t *testing.T) {
	t.Run("sitemapindex.xml exists", func(t *testing.T) {
		path := "./testdata/sitemapindex.xml"
		idx, err := ReadSitemapIndex(path)

		if err != nil {
			t.Errorf("ReadSitemapIndex() should not return error. result:%v", err)
		}

		if len(idx.Sitemap) != 3 {
			t.Errorf("ReadSitemapIndex() should return Sitemap. result:%d expected:%d", 3, len(idx.Sitemap))
		}
	})

	t.Run("sitemapindex.xml not exists", func(t *testing.T) {
		path := "./testdata/not_exist_sitemapindex.xml"
		idx, err := ReadSitemapIndex(path)

		errText := "file not found ./testdata/not_exist_sitemapindex.xml"
		if err.Error() != errText {
			t.Errorf("ReadSitemapIndex() should not return error. result:%s expected:%s", err.Error(), errText)
		}

		if len(idx.Sitemap) != 0 {
			t.Errorf("ReadSitemapIndex() should not return Sitemap. result:%d expected:%d", 0, len(idx.Sitemap))
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("sitemap.xml exists", func(t *testing.T) {
		data, _ := os.ReadFile("./testdata/sitemap.xml")
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

		if err.Error() != "sitemap.xml is empty" {
			t.Errorf("Parse() should return error. result:%s expected:%s", err.Error(), "sitemap.xml is empty")
		}

		if len(smap.URL) != 0 {
			t.Errorf("Parse() should return Sitemap.URL. result:%d expected:%d", 0, len(smap.URL))
		}
	})
}

func TestParseIndex(t *testing.T) {
	t.Run("sitemapindex.xml exists", func(t *testing.T) {
		data, _ := os.ReadFile("./testdata/sitemapindex.xml")
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

		if err.Error() != "sitemapindex.xml is empty" {
			t.Errorf("ParseIndex() should not return error. result:%s expected:%s", err.Error(), "sitemapindex.xml is empty")
		}

		if len(idx.Sitemap) != 0 {
			t.Errorf("ParseIndex() should return Sitemap. result:%d expected:%d", 0, len(idx.Sitemap))
		}
	})
}
