# testdata

The test data was generated using the queries below. Obviously, if they
are run again the results will differ. In short, it is a query for all
articles in the most populare Cardiology journals, with a publish date
in the past year.

The results stored in the json files are used as mock responses for the
tests.


## count.json

The first stage a search query is to determine the total number of results.
This number is then used to determine the number of sets to split the
results into.

The query to get the count for the test data is below:

https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&rettype=count&reldate=365&datetype=pdat&term=loattrfree%20full%20text%5BFilter%5D%20AND%20(%22Am%20Heart%20J%22%5Bjour%5D%20OR%20%22Am%20J%20Cardiol%22%5Bjour%5D%20OR%20%22Arterioscler%20Thromb%20Vasc%20Biol%22%5Bjour%5D%20OR%20%22Atherosclerosis%22%5Bjour%5D%20OR%20%22Basic%20Res%20Cardiol%22%5Bjour%5D%20OR%20%22Cardiovasc%20Res%22%5Bjour%5D%20OR%20%22Chest%22%5Bjour%5D%20OR%20%22Circulation%22%5Bjour%5D%20OR%20%22Circ%20Arrhythm%20Electrophysiol%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Genet%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Qual%20Outcomes%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22Circ%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Circ%20Res%22%5Bjour%5D%20OR%20%22ESC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Acute%20Cardiovasc%20Care%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Pharmacother%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Qual%20Care%20Clin%20Outcomes%22%5Bjour%5D%20OR%20%22Eur%20J%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20J%20Vasc%20Endovasc%20Surg%22%5Bjour%5D%20OR%20%22Europace%22%5Bjour%5D%20OR%20%22Heart%22%5Bjour%5D%20OR%20%22Heart%20Lung%20Circ%22%5Bjour%5D%20OR%20%22Heart%20Rhythm%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22JACC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22J%20Am%20Coll%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Am%20Heart%20Assoc%22%5Bjour%5D%20OR%20%22J%20Am%20Soc%20Echocardiogr%22%5Bjour%5D%20OR%20%22J%20Card%20Fail%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Electrophysiol%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Magn%20Reson%22%5Bjour%5D%20OR%20%22J%20Heart%20Lung%20Transplant%22%5Bjour%5D%20OR%20%22J%20Hypertens%22%5Bjour%5D%20OR%20%22J%20Mol%20Cell%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Thorac%20Cardiovasc%20Surg%22%5Bjour%5D%20OR%20%22J%20Vasc%20Surg%22%5Bjour%5D%20OR%20%22Nat%20Rev%20Cardiol%22%5Bjour%5D%20OR%20%22Prog%20Cardiovasc%20Dis%22%5Bjour%5D%20OR%20%22Resuscitation%22%5Bjour%5D%20OR%20%22Stroke%22%5Bjour%5D)


## idlistx.json

The second stage of fetching the results is fetching a list of article
ids for each of ther required sets.

The three sets in the test data are contained in the files `idlist1.json`,
`idlist2.json` and `idlist3.json`.

The query to fetch these lists is below. This is run three times with the
`retstart=` query paramater set to 0, 1000 and 2000.

As there are 2070 total results for the search, `idlist1.json` contains
the first one thousand ids, `idlist2.json` contains the second thousand
and `idlist3.json` the final 70 - hence 3 'sets'.

ttps://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&retmax=1000&retstart=1000&reldate=365&datetype=pdat&term=loattrfree%20full%20text%5BFilter%5D%20AND%20(%22Am%20Heart%20J%22%5Bjour%5D%20OR%20%22Am%20J%20Cardiol%22%5Bjour%5D%20OR%20%22Arterioscler%20Thromb%20Vasc%20Biol%22%5Bjour%5D%20OR%20%22Atherosclerosis%22%5Bjour%5D%20OR%20%22Basic%20Res%20Cardiol%22%5Bjour%5D%20OR%20%22Cardiovasc%20Res%22%5Bjour%5D%20OR%20%22Chest%22%5Bjour%5D%20OR%20%22Circulation%22%5Bjour%5D%20OR%20%22Circ%20Arrhythm%20Electrophysiol%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Genet%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Qual%20Outcomes%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22Circ%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Circ%20Res%22%5Bjour%5D%20OR%20%22ESC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Acute%20Cardiovasc%20Care%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Pharmacother%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Qual%20Care%20Clin%20Outcomes%22%5Bjour%5D%20OR%20%22Eur%20J%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20J%20Vasc%20Endovasc%20Surg%22%5Bjour%5D%20OR%20%22Europace%22%5Bjour%5D%20OR%20%22Heart%22%5Bjour%5D%20OR%20%22Heart%20Lung%20Circ%22%5Bjour%5D%20OR%20%22Heart%20Rhythm%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22JACC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22J%20Am%20Coll%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Am%20Heart%20Assoc%22%5Bjour%5D%20OR%20%22J%20Am%20Soc%20Echocardiogr%22%5Bjour%5D%20OR%20%22J%20Card%20Fail%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Electrophysiol%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Magn%20Reson%22%5Bjour%5D%20OR%20%22J%20Heart%20Lung%20Transplant%22%5Bjour%5D%20OR%20%22J%20Hypertens%22%5Bjour%5D%20OR%20%22J%20Mol%20Cell%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Thorac%20Cardiovasc%20Surg%22%5Bjour%5D%20OR%20%22J%20Vasc%20Surg%22%5Bjour%5D%20OR%20%22Nat%20Rev%20Cardiol%22%5Bjour%5D%20OR%20%22Prog%20Cardiovasc%20Dis%22%5Bjour%5D%20OR%20%22Resuscitation%22%5Bjour%5D%20OR%20%22Stroke%22%5Bjour%5D)


