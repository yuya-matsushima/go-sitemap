package sitemap

import (
	"os"
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

func BenchmarkForceGet(b *testing.B) {
	server := testServer()
	defer server.Close()

	b.Run("sitemap.xml", func(b *testing.B) {
		url := server.URL + "/sitemap.xml"

		for i := 0; i < b.N; i++ {
			_, err := ForceGet(url, nil)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("contains_empty_sitemap_sitemapindex.xml", func(b *testing.B) {
		url := server.URL + "/contains_empty_sitemap_sitemapindex.xml"

		for i := 0; i < b.N; i++ {
			_, err := ForceGet(url, nil)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkReadSitemap(b *testing.B) {
	path := "./testdata/sitemap.xml"

	for i := 0; i < b.N; i++ {
		_, err := ReadSitemap(path)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkReadSitemapIndex(b *testing.B) {
	path := "./testdata/sitemapindex.xml"

	for i := 0; i < b.N; i++ {
		_, err := ReadSitemapIndex(path)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseSitemap(b *testing.B) {
	data, _ := os.ReadFile("./testdata/sitemap.xml")

	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseSitemapIndex(b *testing.B) {
	data, _ := os.ReadFile("./testdata/sitemapindex.xml")

	for i := 0; i < b.N; i++ {
		_, err := ParseIndex(data)
		if err != nil {
			b.Error(err)
		}
	}
}
