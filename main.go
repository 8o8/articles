package main

import (
	"fmt"
	"github.com/34South/envr"
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/mikedonnici/pubmed"
	"github.com/mikedonnici/elastic"
	"os"
	"strconv"
	"github.com/pkg/errors"
)

// Index maps to the JSON data specifying the index being created
type Index struct {
	Category string `json:"category"`
	Term     string `json:"term"`
	RelDate  int    `json:"reldate"`
}

type Article struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Categories  []string `json:"categories"`
	PubTime     int64    `json:"pubTime"`
	PubDate     string   `json:"pubDate"`
	PubName     string   `json:"pubName"`
	PubNameAbbr string   `json:"pubNameAbbr"`
	PubPageRef  string   `json:"pubPageRef"`
}

func init() {
	envr.New("articlesEnv", []string{
		"ELASTIC_URL",
		"ELASTIC_USER",
		"ELASTIC_PASS",
	}).Auto()
}

func main() {

	fmt.Println("Articler...")

	e := elastic.NewClient(os.Getenv("ELASTIC_URL"), os.Getenv("ELASTIC_USER"), os.Getenv("ELASTIC_PASS"))
	err := e.CheckOK()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to elastic search")

	// Read in the config
	xb, err := ioutil.ReadFile("indices.json")
	if err != nil {
		log.Fatalln("Could not read indices.json")
	}

	var indices []Index
	err = json.Unmarshal(xb, &indices)
	if err != nil {
		log.Fatalln("Could not unmarshal indices.json", err)
	}

	for _, v := range indices {

		p := pubmed.NewSearch(v.Term)
		p.BackDays = v.RelDate
		err := p.Search()
		if err != nil {
			log.Fatalln(err)
		}

		// post batches to elastic
		batchSize := 500
		for i := 0; i < p.ResultCount; i++ {
			xa, err := p.Articles(i, batchSize)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("############################################################################################")
			fmt.Println("Creating Batch", i, "-", i+batchSize)
			body := ""

			for _, a := range xa.Articles {

				doc, err := mapArticle(a)
				if err != nil {
					continue
				}

				// Index the doc
				body += fmt.Sprintf("{\"index\": {\"_id\": \"%v\"}}\n", a.ID)
				body += fmt.Sprintf("%s\n", doc)

				// Update categories in same doc
				script := fmt.Sprintf("{\"source\": \"ctx._source.categories.add(params.category)\", \"lang\": \"painless\", \"params\": {\"category\": \"%s\"}}", v.Category)
				body += fmt.Sprintf("{\"update\": {\"_id\": \"%v\"}}\n", a.ID)
				body += fmt.Sprintf("{\"script\": %s}\n", script)
				body += "\n"

				//fmt.Println(body)
				//os.Exit(0)
			}

			_, err = e.Batch("articles", body)
			if err != nil {
				log.Fatalln(err)
			}

			i += batchSize
		}
	}

	//indices, err := e.Indices()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(indices)
	//
	////err = e.CreateIndex("NewIndex")
	////if err != nil {
	////	fmt.Println(err)
	////	os.Exit(1)
	////}
	//
	//indices, err = e.Indices()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(indices)
	//
	//newDoc := `{"title": "this is the way we roll", "description": "A short story"}`
	//err = e.IndexDoc("articles", "123abc", newDoc)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//xb, err := e.QueryDoc("articles", "123abc")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(string(xb))
	//
	////err = e.DeleteIndex("newindex")
	////if err != nil {
	////	fmt.Println(err)
	////	os.Exit(1)
	////}
	//
	//err = e.UpdateDoc("articles", "123abc", `{"author": "Mike Donnici"}`)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//xb, err = e.QueryDoc("articles", "123abc")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(string(xb))
	//
	//err = e.DeleteDoc("articles", "123abc")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//xb, err = e.QueryDoc("articles", "123abc")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(string(xb))

}

// map a pubmed.Article to local Article, then returns it as a JSON string
func mapArticle(a pubmed.Article) (string, error) {

	var at Article

	at.ID = strconv.Itoa(a.ID)
	at.Title = a.Title
	at.URL = a.URL
	at.Keywords = a.Keywords
	at.PubName = a.Journal
	at.PubNameAbbr = a.JournalAbbrev
	at.PubPageRef = a.Pages
	at.PubTime = a.PubDate.Unix()
	at.PubDate = a.PubDate.String()

	if len(a.Abstract) > 0 {
		at.Description = a.Abstract[0].Value
	}

	// Need this to be empty array, not null, else updating category won't work
	at.Categories = []string{}

	xb, err := json.Marshal(at)
	if err != nil {
		return "", errors.Wrap(err, "mapArticle")
	}

	return string(xb), nil
}
