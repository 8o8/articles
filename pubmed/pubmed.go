package pubmed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"strings"
	"net/url"

	"github.com/pkg/errors"
)

// SearchURL and query parameters
const SearchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&usehistory=y"
const queryReturnMax = "&retmax=%v"
const queryStartIndex = "&retstart=%v"
const queryBackDays = "&reldate=%v&datetype=pdat"
const querySearchTerm = "&term=%v"

// FetchURL fetches articles
const FetchURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&retmode=xml&rettype=abstract&query_key=1&WebEnv=NCID_1_94625809_130.14.18.34_9001_1527039838_862671551_0MetA0_S_MegaStore"

const defaultBackDays = 7
const defaultMaxSetSize = 1000

// Search represents a request to the Pubmed esearch endpoint and results in one or more sets of article IDs which
// can be subsequently used to fetch article summaries.
type Search struct {
	BackDays    int
	Term        string
	ResultCount int
	QueryKey    string `json:"querykey"`
	WebEnv      string `json:"webenv"`
}

// PubMedSummary result - each result is indexed by the id of the record requested, even if there is only one.
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
		BackDays: defaultBackDays,
		Term:     query,
	}
}

//
func (ps *Search) StoreResults() error {

	sURL := SearchURL +
		fmt.Sprintf(queryBackDays, ps.BackDays) +
		fmt.Sprintf(querySearchTerm, ps.Term)

	xb, err := ResponseGET(sURL)
	if err != nil {
		return errors.Wrap(err, "StoreResults")
	}

	var r = struct {
		Result struct {
			Count    string `json:"count"`
			QueryKey string `json:"querykey"`
			WebEnv   string `json:"webenv"`
		} `json:"esearchresult"`
	}{}

	err = json.Unmarshal(xb, &r)
	if err != nil {
		return errors.Wrap(err, "StoreResults, Unmarshal")
	}

	count, err := strconv.Atoi(r.Result.Count)
	if err != nil {
		return errors.Wrap(err, "StoreResults, Atoi")
	}
	ps.ResultCount = count

	ps.QueryKey = r.Result.QueryKey
	ps.WebEnv = r.Result.WebEnv

	return nil
}

// QueryArticles fetches a set of Pubmed article summaries and unmarshals them into a PubMedSet.
// There is no limited to the number of article that may be requested however it is recommended to use the POST method
// for requests of more than 200 articles. The list of article ids is posted as form data in the body of the request.
// Ref: https://www.ncbi.nlm.nih.gov/books/NBK25499/#_chapter4_EFetch_
func (ps *Search) QueryArticles(ids ...string) ([]PubMedArticle, error) {

	var xpa []PubMedArticle

	form := url.Values{}
	form.Add("ids", strings.Join(ids, ","))

	xb, err := ResponsePOST(FetchURL, form)
	if err != nil {
		return xpa, errors.Wrap(err, "QueryArticles could not get response body")
	}

	fmt.Println(string(xb))

	//resultSet.IDs, err = ResultsSetIDs(xb)
	//if err != nil {
	//	return resultSet, errors.Wrap(err, "QueryArticles could not extract count from response body")
	//}

	return nil, nil
}

// ResponseGET returns the response body from a GET request
func ResponseGET(url string) ([]byte, error) {

	httpClient := &http.Client{Timeout: 90 * time.Second}

	r, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "ResponseGET")
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

// ResponsePOST returns the response body from a POST request
func ResponsePOST(url string, data url.Values) ([]byte, error) {

	httpClient := &http.Client{Timeout: 90 * time.Second}

	r, err := httpClient.PostForm(url, data)
	if err != nil {
		return nil, errors.Wrap(err, "ResponseGET")
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

// ResultsSetIDs extracts the list of ids from the JSON response body
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
//func (ps *Search) NumSets() int {
//	return int(math.Ceil(float64(ps.Result.Total) / float64(ps.Result.MaxSetSize)))
//}
