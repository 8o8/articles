package pubmed_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/8o8/articles/pubmed"
	"github.com/matryer/is"
)

const query = `asthma`

var mockResponseJSON = map[string][]byte{
	"count":   []byte{},
	"idlist1": []byte{},
	"idlist2": []byte{},
	"idlist3": []byte{},
}

var mockResponseXML = map[string][]byte{
	"articles": []byte{},
}

func init() {

	for i := range mockResponseJSON {
		f := i + ".json"
		xb, err := ioutil.ReadFile("testdata/" + f)
		if err != nil {
			log.Fatalf("Cannot read %s\n", f)
		}
		mockResponseJSON[i] = xb
	}

	for i := range mockResponseXML {
		f := i + ".xml"
		xb, err := ioutil.ReadFile("testdata/" + f)
		if err != nil {
			log.Fatalf("Cannot read %s\n", f)
		}
		mockResponseXML[i] = xb
	}
}

//func TestResultsSetIDs(t *testing.T) {
//	var err error
//	is := is.New(t)
//	ps := pubmed.NewSearch(query)
//
//	ps.Result.MaxSetSize = 1000
//	ps.Result.Total, err = pubmed.ResultsCount(mockResponseJSON["count"])
//	is.NoErr(err)             // Error setting Total
//	is.Equal(ps.NumSets(), 3) // Expect 3 sets
//
//	xs, err := pubmed.ResultsSetIDs(mockResponseJSON["idlist1"])
//	is.NoErr(err) // Error getting ids from list1
//	is.Equal(len(xs), 1000) // List 1 should have 1000 IDs
//
//	xs, err = pubmed.ResultsSetIDs(mockResponseJSON["idlist2"])
//	is.NoErr(err) // Error getting ids from list1
//	is.Equal(len(xs), 1000) // List 2 should have 1000 IDs
//
//	xs, err = pubmed.ResultsSetIDs(mockResponseJSON["idlist3"])
//	is.NoErr(err) // Error getting ids from list3
//	is.Equal(len(xs), 70) // List 3 should have 70 IDs
//}
//
//
//// Real queries
//
func TestRealSearch(t *testing.T) {
	is := is.New(t)
	ps := pubmed.NewSearch(query)
	ps.BackDays = 100
	err := ps.Search()
	is.NoErr(err)               // Search error
	is.True(ps.ResultCount > 0) // No results for last 100 days?
	set, err := ps.Articles(0, 1000)
	is.NoErr(err) // Fetch error

	for _, a := range set.Articles {
		if len(a.MeshHeadings) > 0 {
			a.Print()
		}
	}
}

//
//func TestFetchArticles(t *testing.T) {
//	is := is.New(t)
//	ps := pubmed.NewSearch(query)
//	ps.Name = "Cardiology"
//	xps, err := ps.Articles("29735362", "29730991")
//	is.NoErr(err) // Error fetching articles
//	fmt.Println(xps)
//}
