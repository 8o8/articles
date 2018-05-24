// Package pubmed fetches articles from Pubmed
package pubmed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const BaseURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils"
const SearchURL = BaseURL + "/esearch.fcgi?db=pubmed&retmode=json&usehistory=y"
const FetchURL = BaseURL + "/efetch.fcgi?db=pubmed&retmode=xml&rettype=abstract"
const queryBackDays = "&reldate=%v&datetype=pdat"
const querySearchTerm = "&term=%v"
const queryReturnMax = "&retmax=%v"
const queryStartIndex = "&retstart=%v"
const queryKey = "&query_key=%v"
const queryWebEnv = "&WebEnv=%s"

const defaultBackDays = 7

// Query represents a search request to the Pubmed esearch endpoint, the results of which are stored at Pubmed
// and subsequently referenced by Key and WebEnv
type Query struct {
	BackDays    int
	Term        string
	ResultCount int
	Key         string `json:"querykey"`
	WebEnv      string `json:"webenv"`
}

// NewSearch returns a pointer to a Query with some defaults set
func NewSearch(query string) *Query {
	return &Query{
		BackDays: defaultBackDays,
		Term:     query,
	}
}

// Search executes the query, results are stored at Pubmed and referenced by the Key and WenEnv values
func (ps *Query) Search() error {

	qURL := SearchURL + fmt.Sprintf(queryBackDays, ps.BackDays) + fmt.Sprintf(querySearchTerm, ps.Term)
	xb, err := responseBody(qURL)
	if err != nil {
		return errors.Wrap(err, "Search")
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
		return errors.Wrap(err, "Search, Unmarshal")
	}

	count, err := strconv.Atoi(r.Result.Count)
	if err != nil {
		return errors.Wrap(err, "Search, Atoi")
	}
	ps.ResultCount = count
	ps.Key = r.Result.QueryKey
	ps.WebEnv = r.Result.WebEnv

	return nil
}

// Articles fetches a set of articles from the Pubmed cache referenced by Key and WebEnv. The response is an xml
// payload that is unmarshaled into a PubMedSet.
// Ref: https://www.ncbi.nlm.nih.gov/books/NBK25499/#_chapter4_EFetch_
func (ps *Query) Articles(startIndex, retMax int) (ArticleSet, error) {

	var set ArticleSet

	qURL := FetchURL +
		fmt.Sprintf(queryKey, ps.Key) +
		fmt.Sprintf(queryWebEnv, ps.WebEnv) +
		fmt.Sprintf(queryStartIndex, startIndex) +
		fmt.Sprintf(queryReturnMax, retMax)
	xb, err := responseBody(qURL)
	if err != nil {
		return set, errors.Wrap(err, "Articles could not get response body")
	}

	err = xml.Unmarshal(xb, &set)
	if err != nil {
		return set, errors.Wrap(err, "Articles could not unmarshal response body")
	}

	for i, a := range set.Articles {
		set.Articles[i].PubDate, _ = bestPubDate(a) // todo ... ignore the error?
		set.Articles[i].Keywords = mergeKeywords(a)
		set.Articles[i].URL = articleURL(a)
		set.Articles[i].Citation = citation(a)
	}

	return set, nil
}

// responseBody returns the response body from a GET request as a []byte
func responseBody(url string) ([]byte, error) {

	httpClient := &http.Client{Timeout: 90 * time.Second}
	r, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "responseBody")
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}
