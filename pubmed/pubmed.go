package pubmed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// SearchURL and query parameters
const SearchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json"
const queryReturnMax = "&retmax=%v"
const queryStartIndex = "&retstart=%v"
const queryBackDays = "&reldate=%v&datetype=pdat"
const querySearchTerm = "&term=%v"

// FetchURL is the endpoint that fetches a single pubmed article by id
const FetchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&retmode=xml&rettype=abstract&id="

const defaultSearchName = "Pubmed Search"
const defaultBackDays = 7
const defaultMaxSetSize = 1000

// Search represents the Pubmed query
type Search struct {
	Name     string
	BackDays int
	Term     string
	Result   Result
}

// Result of the Search where MaxSetSize determines the maximum number of IDs that can be stored in each set.
type Result struct {
	Total      int
	MaxSetSize int
	Sets       []Set
}

// Set is a set of article IDs returns from the Search
type Set struct {
	IDs []string `json:"idlist"`
}

// PubMedSummary Result - each result is indexed by the id of the record requested, even if there is only one.
// As we can pass multiple ids on the URL to save requests eg 521345,765663,121234,3124256
// use XMl to fetch this as easier to get to nested data
type PubMedArticleSet struct {
	Articles []PubMedArticle `xml:"PubmedArticle"`
}
type PubMedArticle struct {
	ID              int            `xml:"MedlineCitation>PMID"`
	Title           string         `xml:"MedlineCitation>Article>ArticleTitle"`
	Abstract        []AbstractText `xml:"MedlineCitation>Article>Abstract>AbstractText"`
	ArticleIDList   []ArticleID    `xml:"PubmedData>ArticleIdList>ArticleId"`
	KeywordList     []string       `xml:"MedlineCitation>KeywordList>Keyword"`
	MeshHeadingList []string       `xml:"MedlineCitation>MeshHeadingList>MeshHeading>DescriptorName"`
	AuthorList      []Author       `xml:"MedlineCitation>Article>AuthorList>Author"`
	Journal         string         `xml:"MedlineCitation>Article>Journal>Title"`
	JournalAbbrev   string         `xml:"MedlineCitation>Article>Journal>ISOAbbreviation"`
	Volume          string         `xml:"MedlineCitation>Article>Journal>JournalIssue>Volume"`
	Issue           string         `xml:"MedlineCitation>Article>Journal>JournalIssue>Issue"`
	Pages           string         `xml:"MedlineCitation>Article>Pagination>MedlinePgn"`
	PubYear         string         `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Year"`
	PubMonth        string         `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Month"`
	PubDay          string         `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Day"`

	// Note that for these fallback dates there are multiple nodes as each is part of the records history
	// Ideally, we would pick the xml node with the attribute 'entrez', which is the oldest.
	// today loo into ensuring the oldest date element is selected for fallback
	PubYearFallback  string `xml:"PubmedData>History>PubMedPubDate>Year"`
	PubMonthFallback string `xml:"PubmedData>History>PubMedPubDate>Month"`
	PubDayFallback   string `xml:"PubmedData>History>PubMedPubDate>Day"`
}

type AbstractText struct {
	Key   string `xml:"label,attr"`
	Value string `xml:",chardata"`
}

type ArticleID struct {
	Key   string `xml:"IdType,attr"`
	Value string `xml:",chardata"`
}

type Author struct {
	Key      string `xml:"ValidYN,attr"`
	LastName string `xml:"LastName"`
	Initials string `xml:"Initials"`
}

type Resource struct {
	Published   time.Time              `json:"published" bson:"published"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Keywords    []string               `json:"keywords"`
	URL         string                 `json:"url"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// NewSearch returns a pointer to a Search with some defaults set
func NewSearch(query string) *Search {
	return &Search{
		Name:     defaultSearchName,
		BackDays: defaultBackDays,
		Term:     query,
		Result: Result{
			MaxSetSize: defaultMaxSetSize,
		},
	}
}

// Query runs the pubmed query
func (ps *Search) Query() error {

	err := ps.QueryResultTotal()
	if err != nil {
		return err
	}

	c := ps.NumSets()
	for i := 0; i < c; i++ {
		set, err := ps.QueryResultSet(i)
		if err != nil {
			return err
		}
		ps.Result.Sets = append(ps.Result.Sets, set)
	}

	for i, v := range ps.Result.Sets {
		fmt.Println("Set #", i, len(v.IDs))
	}

	return nil
}

//
func (ps *Search) QueryResultTotal() error {

	url := SearchURL + fmt.Sprintf(queryBackDays+querySearchTerm, ps.BackDays, ps.Term) + "&rettype=count"

	xb, err := ResponseBody(url)
	if err != nil {
		return errors.Wrap(err, "QueryResultTotal could not get response body")
	}

	ps.Result.Total, err = ResultsCount(xb)
	if err != nil {
		return errors.Wrap(err, "QueryResultTotal could not extract count from response body")
	}

	return nil
}

//
func (ps *Search) QueryResultSet(setIndex int) (Set, error) {

	var resultSet Set

	startIndex := setIndex * ps.Result.MaxSetSize

	url := SearchURL +
		fmt.Sprintf(queryBackDays, ps.BackDays) +
		fmt.Sprintf(queryStartIndex, startIndex) +
		fmt.Sprintf(queryReturnMax, ps.Result.MaxSetSize) +
		fmt.Sprintf(querySearchTerm, ps.Term)

	xb, err := ResponseBody(url)
	if err != nil {
		return resultSet, errors.Wrap(err, "QueryResultSet could not get response body")
	}

	resultSet.IDs, err = ResultsSetIDs(xb)
	if err != nil {
		return resultSet, errors.Wrap(err, "QueryResultTotal could not extract count from response body")
	}

	return resultSet, nil
}

func ResponseBody(url string) ([]byte, error) {

	httpClient := &http.Client{Timeout: 90 * time.Second}

	r, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "ResponseBody")
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

// extract the count value from response body
func ResultsCount(responseBody []byte) (int, error) {

	var r = struct {
		Result struct {
			Count string `json:"count"`
		} `json:"esearchresult"`
	}{}

	err := json.Unmarshal(responseBody, &r)
	if err != nil {
		return 0, errors.Wrap(err, "Unmarshal")
	}

	count, err := strconv.Atoi(r.Result.Count)
	if err != nil {
		return 0, errors.Wrap(err, "Atoi")
	}

	return count, nil
}

// extract the sets from response body
func ResultsSetIDs(responseBody []byte) ([]string, error) {

	var r = struct {
		Result struct {
			IDList []string `json:"idlist"`
		} `json:"esearchresult"`
	}{}

	err := json.Unmarshal(responseBody, &r)
	if err != nil {
		return []string{}, errors.Wrap(err, "Unmarshal")
	}

	return r.Result.IDList, nil
}

// NumSets calculates the number of sets required for a Search based on MaxSetSize and total results.
func (ps *Search) NumSets() int {
	return int(math.Ceil(float64(ps.Result.Total) / float64(ps.Result.MaxSetSize)))
}
