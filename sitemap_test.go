package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSitemap(t *testing.T) {
	server := server()
	defer server.Close()

	data, err := GetSitemap(server.URL + "/sitemap.xml")

	if len(data.URL) == 0 {
		t.Error("GetSitemap() should return Some Sitemap.Url data")
	}

	if err != nil {
		t.Error("GetSitemap() should not has error")
	}
}

func TestGetSitemapRecivedInvalidSitemapURL(t *testing.T) {
	server := server()
	defer server.Close()

	_, err := GetSitemap(server.URL + "/emptymap.xml")

	if err == nil {
		t.Error("GetSitemap() should return error")
	}
}

func TestGetSitemapRecivedSitemapIndexURL(t *testing.T) {
	server := server()
	defer server.Close()

	data, err := GetSitemap(server.URL + "/sitemapindex.xml")

	if len(data.URL) == 0 {
		t.Error("GetSitemap() should return Some Sitemap.Url data")
	}

	if err != nil {
		t.Error("GetSitemap() should not has error")
	}
}

func BenchmarkReadSitemapXML(b *testing.B) {
	server := server()
	defer server.Close()
	GetSitemap(server.URL + "/sitemap.xml")
}

func BenchmarkReadSitemapIndex(b *testing.B) {
	server := server()
	defer server.Close()
	GetSitemap(server.URL + "/sitemapindex.xml")
}

func server() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "" {
			// index page is always not found
			http.NotFound(w, r)
		}

		res, err := ioutil.ReadFile("./testdata" + r.RequestURI)
		if err != nil {
			http.NotFound(w, r)
		}
		str := strings.Replace(string(res), "HOST", r.Host, -1)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, str)
	}))

	return server
}
