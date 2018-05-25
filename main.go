package main

import (
	"fmt"
	"github.com/mikedonnici/pubmed"
)

func main() {
	fmt.Println("Fetching articles from pubmed...")

	p := pubmed.NewSearch("shoulder%20subluxation")
	p.BackDays = 1000
	p.Search()
	xa, err := p.Articles(0, 1)
	if err != nil {
		fmt.Println(err)
	}

	for _, a := range xa.Articles {
		a.Print()
	}
}
