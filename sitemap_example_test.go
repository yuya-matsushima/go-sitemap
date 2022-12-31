package sitemap

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func ExampleGet() {
	smap, err := Get("https://issueoverflow.com/sitemap.xml", nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, URL := range smap.URL {
		fmt.Println(URL.Loc)
	}
}

func ExampleGet_changeFetch() {
	SetFetch(func(URL string, options interface{}) ([]byte, error) {
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return []byte{}, err
		}

		// Set User-Agent
		req.Header.Set("User-Agent", "MyBot")

		// Set timeout
		timeout := time.Duration(10 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		// Fetch data
		res, err := client.Do(req)
		if err != nil {
			return []byte{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return []byte{}, err
		}

		return body, err
	})

	smap, err := Get("https://issueoverflow.com/sitemap.xml", nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, URL := range smap.URL {
		fmt.Println(URL.Loc)
	}
}
