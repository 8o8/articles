# pubmed

## Pubmed Queries

The Pubmed search utilities are powerful but can be a bit curlt to use.

After a lot of trial and error the following approach was adopted.

Initially, Pubmed requests were being made in batches, then the results
stored in batches, iterated over and stored in an index.

Pubmed can actually store the entire results set for you and you can
simply request subsets of those results and process them.

For example, the query below searched for asthma articles with a publish date
within the past 7 days:

```http
https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&datetype=pdat&reldate=7&retstart=0&retmax=20&term=asthma
```

**params**

* `datetype=pdate` - query publish date
* `reldate=7` - publish date within last 7 days
* `retstart=0` - return set start at first record
* `retmax=20` - return a maximum of 20 results

This returns a maximum of 20 article IDs, starting with the first result,
as below:

```json
{
    "header": {...},
    "esearchresult": {
        "count": "90",
        "retmax": "20",
        "retstart": "0",
        "idlist": ["29783109", "29782937", ... "8 more"]
    }
}
```

This set of IDs was stored, and the next set of IDs fetched and stored,
and so on until we have all 90 IDs retrieved.

Thereafter the actual article summaries were retrieved, also in batches,
and indexed.

Although this works well, it requires batches to be handled at the
request side and is a tad cumbersome.

A better way...

The url can be simplified a bit by removing `retstart` and `retmax` -
they are irrelevant for the next step and will just default to 0 and 20
respectively.

A new param is added to the request url - `userhistory=y`. This will tell
 Pubmed to store the results for subsequent retrieval.

```http
https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&datetype=pdat&reldate=7&usehistory=y&term=asthma
```

The response contains two additional fields, **`querykey`** and **`webenv`**.
These fields are used to reference the stored results.

```json
{
    "header": {...},
    "esearchresult": {
        "count": "90",
        "retmax": "20",
        "retstart": "0",
        "querykey": "1",
        "webenv": "NCID_1_95852708_130.14.18.34_9001_1527045975_1203341711_0MetA0_S_MegaStore",
        "idlist": ["29783109", "29782937", ... "8 more"]
        ]
    }
}
```

**Important** - although the `idlist` only displays the 20 results limited
by `retmax`, the stored search results contain *all* 90 of articles found.

The stored search results can then be retrieved by adding the url params:
 * `query_key=1`
 * `WebEnv=NCID_1_95852708_130.14.18.34_9001_1527045975_1203341711_0MetA0_S_MegaStore`

```http
https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.cgi?db=pubmed&retmode=xml&rettype=abstract&query_key=1&WebEnv=NCID_1_95852708_130.14.18.34_9001_1527045975_1203341711_0MetA0_S_MegaStore
```

...and voila! All 90 results are there, regardless of `retmax` or `retstart`
values in the previous step.

**Caution** - If the initial search returned a large result set then this
fetch can return a substantial file.

To fetch subsets of the articles in the result set use `retmax` and `retstart`
  in the same way as above.

The first 10 results:

* `retstart=0`
* `retmax=10`

... the next 10:
* `retstart=10`
* `retmax=10`

.. and so on.

```http
https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.cgi?db=pubmed&retmode=xml&rettype=abstract&retstart=0&retmax=10&query_key=1&WebEnv=NCID_1_95852708_130.14.18.34_9001_1527045975_1203341711_0MetA0_S_MegaStore
```

So some batchin processing will still be required, meaning that knowing
the total number of results is important, however it is nowehere as
 cumbersome as searching *and* fetching in batches.




