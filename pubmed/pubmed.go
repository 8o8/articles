package pubmed

import (
	"time"
	"log"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SearchURL and query parameters
const SearchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json"
const queryReturnMax = "&retmax=%v"
const queryStartIndex = "&retstart=%v"
const queryBackDays = "&reldate=%v&datetype=pdat"
const querySearchTerm = "&term=%v"

// FetchURL is the endpoint that fetches a single pubmed article by id
const FetchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&retmode=xml&rettype=abstract&id="

var httpClient = &http.Client{Timeout: 90 * time.Second}

// Search represents the search being done at pubmed
type Search struct {
	Category string
	Term     string
	BackDays int
	Count    int
	Response response
}

// response maps the fields in the response from pubmed
type response struct {
	Header map[string]string `json:"header"`
	Result result            `json:"esearchresult"`
}

// result maps the esearchresult field in the response
type result struct {
	Count string   `json:"count"`
	IDs   []string `json:"idlist"`
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

func NewSearch() *Search {
	return &Search{}
}

// SetCount runs the pubmed query with rettype=count to get the total article count for the search term.
func (ps *Search) SetCount() error {

	var c = struct {
		Result map[string]string `json:"esearchresult"`
	}{}

	url := SearchURL + fmt.Sprintf(queryBackDays+querySearchTerm, ps.BackDays, ps.Term) + "&rettype=count"
	r, err := httpClient.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	count, err := strconv.Atoi(c.Result["count"])
	if err != nil {
		log.Fatalln(err)
	}

	ps.Count = count
	return nil
}
