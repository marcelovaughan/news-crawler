package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func crawlerAction(w http.ResponseWriter, r *http.Request) {

	// //Allow CORS here By * or specific origin
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var links []Link

	c := colly.NewCollector(
		//colly.AllowedDomains("www.globo.com", "g1.globo.com"),
		colly.MaxDepth(2),
	)

	keyword := r.FormValue("keyword")
	site := r.FormValue("site")

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 10
	}

	if len(keyword) < 1 {
		fmt.Println("Please specify a keyword and a website")
		//os.Exit(1)
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Contains(strings.ToUpper(e.Text), strings.ToUpper(keyword)) {
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			fmt.Printf("Links: %b\n", len(links))

			if len(links) < limit {
				var linkFound Link
				linkFound.Title = e.Text
				linkFound.URL = link

				links = append(links, linkFound)
				//c.Visit(e.Request.AbsoluteURL(link))
			} else {
				return
			}

		}

	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit(site)

	// dump results
	b, err := json.Marshal(links)
	if err != nil {
		//log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func crawler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!!!!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", crawler)
	router.HandleFunc("/news", crawlerAction).Methods("POST")

	var port = "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
