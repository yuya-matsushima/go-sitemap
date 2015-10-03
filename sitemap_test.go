package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	server := server()
	defer server.Close()

	data, err := Get(server.URL + "/sitemap.xml")

	if len(data.URL) == 0 {
		t.Error("Get() should return Some Sitemap.Url data")
	}

	if err != nil {
		t.Error("Get() should not has error")
	}
}

func TestGetRecivedInvalidSitemapURL(t *testing.T) {
	server := server()
	defer server.Close()

	_, err := Get(server.URL + "/emptymap.xml")

	if err == nil {
		t.Error("Get() should return error")
	}
}

func TestGetRecivedSitemapIndexURL(t *testing.T) {
	server := server()
	defer server.Close()

	data, err := Get(server.URL + "/sitemapindex.xml")

	if len(data.URL) == 0 {
		t.Error("Get() should return Some Sitemap.Url data")
	}

	if err != nil {
		t.Error("Get() should not has error")
	}
}

func BenchmarkReadSitemapXML(b *testing.B) {
	server := server()
	defer server.Close()
	Get(server.URL + "/sitemap.xml")
}

func BenchmarkReadSitemapIndex(b *testing.B) {
	server := server()
	defer server.Close()
	Get(server.URL + "/sitemapindex.xml")
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
