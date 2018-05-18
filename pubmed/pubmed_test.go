package pubmed_test

import (
	"fmt"
	"testing"
	"io/ioutil"
	"log"

	"github.com/8o8/articles/pubmed"
	"github.com/matryer/is"
)

const query = `loattrfree%20full%20text%5BFilter%5D%20AND%20(%22Am%20Heart%20J%22%5Bjour%5D%20OR%20%22Am%20J%20Cardiol%22%5Bjour%5D%20OR%20%22Arterioscler%20Thromb%20Vasc%20Biol%22%5Bjour%5D%20OR%20%22Atherosclerosis%22%5Bjour%5D%20OR%20%22Basic%20Res%20Cardiol%22%5Bjour%5D%20OR%20%22Cardiovasc%20Res%22%5Bjour%5D%20OR%20%22Chest%22%5Bjour%5D%20OR%20%22Circulation%22%5Bjour%5D%20OR%20%22Circ%20Arrhythm%20Electrophysiol%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Genet%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Qual%20Outcomes%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22Circ%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Circ%20Res%22%5Bjour%5D%20OR%20%22ESC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Acute%20Cardiovasc%20Care%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Pharmacother%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Qual%20Care%20Clin%20Outcomes%22%5Bjour%5D%20OR%20%22Eur%20J%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20J%20Vasc%20Endovasc%20Surg%22%5Bjour%5D%20OR%20%22Europace%22%5Bjour%5D%20OR%20%22Heart%22%5Bjour%5D%20OR%20%22Heart%20Lung%20Circ%22%5Bjour%5D%20OR%20%22Heart%20Rhythm%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22JACC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22J%20Am%20Coll%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Am%20Heart%20Assoc%22%5Bjour%5D%20OR%20%22J%20Am%20Soc%20Echocardiogr%22%5Bjour%5D%20OR%20%22J%20Card%20Fail%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Electrophysiol%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Magn%20Reson%22%5Bjour%5D%20OR%20%22J%20Heart%20Lung%20Transplant%22%5Bjour%5D%20OR%20%22J%20Hypertens%22%5Bjour%5D%20OR%20%22J%20Mol%20Cell%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Thorac%20Cardiovasc%20Surg%22%5Bjour%5D%20OR%20%22J%20Vasc%20Surg%22%5Bjour%5D%20OR%20%22Nat%20Rev%20Cardiol%22%5Bjour%5D%20OR%20%22Prog%20Cardiovasc%20Dis%22%5Bjour%5D%20OR%20%22Resuscitation%22%5Bjour%5D%20OR%20%22Stroke%22%5Bjour%5D)`

var mockResponses = map[string][]byte{
	"count":   []byte{},
	"idlist1": []byte{},
	"idlist2": []byte{},
	"idlist3": []byte{},
}

func init() {
	for i := range mockResponses {
		f := i + ".json"
		xb, err := ioutil.ReadFile("testdata/" + f)
		if err != nil {
			log.Fatalf("Cannot read %s\n", f)
		}
		mockResponses[i] = xb
	}
}

func TestSetCount(t *testing.T) {
	is := is.New(t)
	ps := pubmed.NewSearch(query)
	ps.Name = "Cardiology"
	ps.BackDays = 100
	err := ps.Query()
	is.NoErr(err)                // Error running query
	is.True(ps.Result.Total > 0) // No results for last 100 days?
	fmt.Println(ps.Result.Total)
}

func TestNumberOfSets(t *testing.T) {
	var err error
	is := is.New(t)
	ps := pubmed.NewSearch(query)

	ps.Result.MaxSetSize = 1000
	ps.Result.Total, err = pubmed.ResultsCount(mockResponses["count"])
	is.NoErr(err)             // Error setting Total
	is.Equal(ps.NumSets(), 3) // Expect 3 sets

	ps.Result.MaxSetSize = 50
	ps.Result.Total, err = pubmed.ResultsCount(mockResponses["count"])
	is.NoErr(err)              // Error setting Total
	is.Equal(ps.NumSets(), 42) // Expect 42 sets
}

func TestResultsSetIDs(t *testing.T) {
	var err error
	is := is.New(t)
	ps := pubmed.NewSearch(query)

	ps.Result.MaxSetSize = 1000
	ps.Result.Total, err = pubmed.ResultsCount(mockResponses["count"])
	is.NoErr(err)             // Error setting Total
	is.Equal(ps.NumSets(), 3) // Expect 3 sets

	xs, err := pubmed.ResultsSetIDs(mockResponses["idlist1"])
	is.NoErr(err) // Error getting ids from list1
	is.Equal(len(xs), 1000) // List 1 should have 1000 IDs

	xs, err = pubmed.ResultsSetIDs(mockResponses["idlist2"])
	is.NoErr(err) // Error getting ids from list1
	is.Equal(len(xs), 1000) // List 2 should have 1000 IDs

	xs, err = pubmed.ResultsSetIDs(mockResponses["idlist3"])
	is.NoErr(err) // Error getting ids from list3
	is.Equal(len(xs), 70) // List 3 should have 70 IDs
}
