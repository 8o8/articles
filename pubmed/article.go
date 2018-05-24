package pubmed

import (
	"fmt"
	"time"
	"strings"
)

// ArticleSet maps to the root node of the returned XML
type ArticleSet struct {
	Articles []Article `xml:"PubmedArticle"`
}

// Article maps on of the <PubmedArticle> nodes
type Article struct {
	ID            int             `xml:"MedlineCitation>PMID"`
	ArticleIDs    []ArticleID     `xml:"PubmedData>ArticleIdList>ArticleId"`
	Title         string          `xml:"MedlineCitation>Article>ArticleTitle"`
	Abstract      []AbstractParts `xml:"MedlineCitation>Article>Abstract>AbstractParts"`
	Keywords      []string        `xml:"MedlineCitation>KeywordList>Keyword"`
	MeshHeadings  []string        `xml:"MedlineCitation>MeshHeadingList>MeshHeading>DescriptorName"`
	Authors       []Author        `xml:"MedlineCitation>Article>AuthorList>Author"`
	Journal       string          `xml:"MedlineCitation>Article>Journal>Title"`
	JournalAbbrev string          `xml:"MedlineCitation>Article>Journal>ISOAbbreviation"`
	Volume        string          `xml:"MedlineCitation>Article>Journal>JournalIssue>Volume"`
	Issue         string          `xml:"MedlineCitation>Article>Journal>JournalIssue>Issue"`
	Pages         string          `xml:"MedlineCitation>Article>Pagination>MedlinePgn"`
	PubYear       string          `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Year"`
	PubMonth      string          `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Month"`
	PubDay        string          `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Day"`

	// Note that for these fallback dates there are multiple nodes as each is part of the records history
	// Ideally, we would pick the xml node with the attribute 'entrez', which is the oldest.
	// today loo into ensuring the oldest date element is selected for fallback
	PubYearFallback  string `xml:"PubmedData>History>PubMedPubDate>Year"`
	PubMonthFallback string `xml:"PubmedData>History>PubMedPubDate>Month"`
	PubDayFallback   string `xml:"PubmedData>History>PubMedPubDate>Day"`

	// These will be set after the fact
	PubDate  time.Time
	Citation string
	URL      string
}

// AbstractParts represents a single section in the article abstract. It may contain one or more sections or
// paragraphs and these are returned as separate xml nodes.
type AbstractParts struct {
	Key   string `xml:"label,attr"`
	Value string `xml:",chardata"`
}

// ArticleID maps to the various types of identifiers associated with an article - "pubmed", "pii", "doi" etc
type ArticleID struct {
	Key   string `xml:"IdType,attr"`
	Value string `xml:",chardata"`
}

// Author is one of the contributors to the article
type Author struct {
	Key      string `xml:"ValidYN,attr"`
	LastName string `xml:"LastName"`
	Initials string `xml:"Initials"`
}

func (a *Article) Print() {
	fmt.Println("==================================================================================================")
	fmt.Println("Pubmed ID:", a.ID)
	fmt.Println("All IDs:", a.ArticleIDs)
	fmt.Println("Title:", a.Title)
	l := len(a.Abstract)
	fmt.Printf("Abstract (%v): %s\n", l, a.Abstract)
	l = len(a.Keywords)
	fmt.Printf("Keywords (%v): %s\n", l, a.Keywords)
	l = len(a.MeshHeadings)
	fmt.Printf("MESH Headings (%v): %s\n", l, a.MeshHeadings)
	l = len(a.Authors)
	fmt.Printf("Authors (%v): %s\n", l, a.Authors)
	fmt.Println("Journal:", a.Journal)
	fmt.Println("Journal Abbrev:", a.JournalAbbrev)
	fmt.Println("Volume:", a.Volume)
	fmt.Println("Issue:", a.Issue)
	fmt.Println("Pages:", a.Pages)
	fmt.Println("Publish Date:", a.PubDate)
	fmt.Println("PubYear:", a.PubYear)
	fmt.Println("PubMonth:", a.PubMonth)
	fmt.Println("PubDay:", a.PubDay)
	fmt.Println("PubYearFallback:", a.PubYearFallback)
	fmt.Println("PubMonthFallback:", a.PubMonthFallback)
	fmt.Println("PubDayFallback:", a.PubDayFallback)
	fmt.Println("Citation (part):", a.Citation)
	fmt.Println("URL:", a.URL)
}

// bestPubDate attempts to find the best date value to set as the publish date. Records with bung date fields won't parse
// so the fallback values might be used. The fallback values are a set of dates in the history of the pubmed article,
// keyed with strings that indicate progression through the publication workflow, eg "entrez" -> "pubmed" -> "medline".
func bestPubDate(a Article) (time.Time, error) {

	var err error

	// Day is often missing from the Pubmed data, set to 1 so we can create time values.
	// However, leave the original value empty for the descriptive PubDate in attributes below.
	day := a.PubDay
	if day == "" {
		day = "1"
	}

	// Month in Pubmed data is *usually* a 3 character string like "May", but can also be a 1 or 2 character
	// string number, "5" or "05". We want a string number to convert to a time value.
	months := map[string]string{
		"Jan": "1", "Feb": "2", "Mar": "3", "Apr": "4", "May": "5", "Jun": "6",
		"Jul": "7", "Aug": "8", "Sep": "9", "Oct": "10", "Nov": "11", "Dec": "12",
	}
	month, ok := months[a.PubMonth]
	if !ok {
		month = a.PubMonth
	}

	// A value for years seems to always be present
	year := a.PubYear

	d := year + "-" + month + "-" + day
	pubDate, err := time.Parse("2006-1-2", d)
	if err == nil {
		return pubDate, nil // success
	}

	// Try fallback dates
	d = a.PubYearFallback + "-" + a.PubMonthFallback + "-" + a.PubDayFallback
	return time.Parse("2006-1-2", d)
}


// mergeKeyWords adds more values to the returned keywords to assist searches. Medline articles return only
// MeshHeadings (no keywords) so these are added in, as well as authors and the various id values.
func mergeKeywords(a Article) []string {

	xs := a.Keywords
	xs = append(xs, a.MeshHeadings...)

	for _, v := range a.Authors {
		xs = append(xs, v.LastName + " " + v.Initials)
	}

	for _, v := range a.ArticleIDs {
		xs = append(xs, v.Value)
	}

	// Some of the mesh terms are phrases that have a comma, eg "Intubation, Intratracheal" - these are removed
	// to minimise any affect on searches or unmarshal operations
	for i, w := range xs {
		xs[i] = strings.TrimSpace(strings.Replace(w, ",", "", -1))
	}

	return xs
}

// articleURL returns the url or DOI link
func articleURL(a Article) string {
	for _, v := range a.ArticleIDs {
		if v.Key == "doi" {
			return "https://doi.org/" + v.Value
		}
	}
	return ""
}

// citation returns a suitably formatted citation string
func citation(a Article) string {
	return fmt.Sprintf("%s. %s; %s(%s): %s", a.Journal, a.PubYear, a.Volume, a.Issue, a.Pages)
}
