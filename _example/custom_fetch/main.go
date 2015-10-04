package main

import (
	"fmt"
	"github.com/yterajima/go-sitemap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	server := server()
	defer server.Close()

	sitemap.SetFetch(myFetch)

	smap, err := sitemap.Get(server.URL + "/sitemap.xml")
	if err != nil {
		fmt.Println(err)
	}

	// Print URL in sitemap.xml
	for _, URL := range smap.URL {
		fmt.Println(URL.Loc)
	}
}

func myFetch(URL string) ([]byte, error) {
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func server() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Print User-Agent
		fmt.Println("User-Agent: " + r.Header.Get("User-Agent"))

		res, err := ioutil.ReadFile("../../testdata" + r.RequestURI)
		if err != nil {
			http.NotFound(w, r)
		}
		str := strings.Replace(string(res), "HOST", r.Host, -1)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, str)
	}))

	return server
}
