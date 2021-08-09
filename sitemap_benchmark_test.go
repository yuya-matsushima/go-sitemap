package sitemap

import (
	"io/ioutil"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	server := testServer()
	defer server.Close()

	b.Run("sitemap.xml", func(b *testing.B) {
		url := server.URL + "/sitemap.xml"

		for i := 0; i < b.N; i++ {
			_, err := Get(url, nil)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("sitemapindex.xml", func(b *testing.B) {
		url := server.URL + "/sitemapindex.xml"

		for i := 0; i < b.N; i++ {
			_, err := Get(url, nil)
			if err != nil {
				b.Error(err)
			}
		}
	})
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
