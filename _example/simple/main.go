package main

import (
	"fmt"

	"github.com/yterajima/go-sitemap"
)

func main() {
	smap, err := sitemap.Get("http://www.e2esound.com/sitemap.xml", nil)
	if err != nil {
		fmt.Println(err)
	}

	// Print URL in sitemap.xml
	for _, URL := range smap.URL {
		fmt.Println(URL.Loc)
	}
}
