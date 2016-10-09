package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

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
