package sitemap

import (
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"strings"
)

func testServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "" {
			// index page is always not found
			http.NotFound(w, r)
			return
		}

		res, err := os.ReadFile("./testdata" + r.RequestURI)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		str := strings.Replace(string(res), "HOST", r.Host, -1)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, str)
	}))
}
