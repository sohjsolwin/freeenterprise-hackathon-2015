package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var stateUnemploymentRequestTemplate string = "https://www.quandl.com/api/v3/datasets/FRED/%sUR.json?auth_token=7GKGd3eZ4wexWY4Ge1bb&start_date=%s&end_date=%s"

var havenNewQueryTemplate string = "https://api.havenondemand.com/1/api/sync/querytextindex/v1?text=%s&ignore_operators=false&indexes=news_eng&max_date=%s&print=fields&print_fields=reference+title&promotion=false&sort=date&summary=off&total_results=false&apikey=93e3f03c-9767-4c2c-b38d-bd48c597295e"

var stateData map[string]([]StateUnemploymentDataRecord)
var stateSentimentData map[string](map[string]StateArticleSentimentRecord)

var stateList = map[string]string{
	"AL": "%22alabama%22",
	"AK": "%22alaska%22",
	"AZ": "%22arizona%22",
	"AR": "%22arkansas%22",
	"CA": "%22california%22",
	"CO": "%22colorado%22",
	"CT": "%22connecticut%22",
	"DE": "%22delaware%22",
	"FL": "%22florida%22",
	"GA": "%22georgia%22",
	"HI": "%22hawaii%22",
	"ID": "%22idaho%22",
	"IL": "%22illinois%22",
	"IN": "%22indiana%22",
	"IA": "%22iowa%22",
	"KS": "%22kansas%22",
	"KY": "%22kentucky%22",
	"LA": "%22louisiana%22",
	"ME": "%22maine%22",
	"MD": "%22maryland%22",
	"MA": "%22massachusetts%22",
	"MI": "%22michigan%22",
	"MN": "%22minnesota%22",
	"MS": "%22mississippi%22",
	"MO": "%22missouri%22",
	"MT": "%22montana%22",
	"NE": "%22nebraska%22",
	"NV": "%22nevada%22",
	"NH": "%22new+hampshire%22",
	"NJ": "%22new+jersey%22",
	"NM": "%22new+mexico%22",
	"NY": "%22new+york%22",
	"NC": "%22north+carolina%22",
	"ND": "%22north+dakota%22",
	"OH": "%22ohio%22",
	"OK": "%22oklahoma%22",
	"OR": "%22oregon%22",
	"PA": "%22pennsylvania%22",
	"RI": "%22rhode+island%22",
	"SC": "%22south+carolina%22",
	"SD": "%22south+dakota%22",
	"TN": "%22tennessee%22",
	"TX": "%22texas%22",
	"UT": "%22utah%22",
	"VT": "%22vermont%22",
	"VA": "%22virginia%22",
	"WA": "%22washington%22",
	"WV": "%22west+virginia%22",
	"WI": "%22wisconsin%22",
	"WY": "%22wyoming%22",
	"DC": "%22washington+dc%22",
}

type StateUnemploymentDataRecord struct {
	Date  string
	Value float64
}
type StateArticleSentimentRecord struct {
	SentimentScore float32
	ArticleLink    string
}

type QuandlEmploymentData struct {
	Dataset QuandlDataset `json:"dataset"`
}

type QuandlDataset struct {
	Id                  int             `json:"id"`
	DatasetCode         string          `json:"dataset_code"`
	DatabaseCode        string          `json:"database_code"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	RefreshedAt         string          `json:"refreshed_at"`
	NewestAvailableDate string          `json:"newest_available_date"`
	OldestAvailableDate string          `json:"oldest_available_date"`
	ColumnNames         []string        `json:"column_names"`
	Frequency           string          `json:"frequency"`
	Type                string          `json:"type"`
	Premium             bool            `json:"premium"`
	Limit               string          `json:"limit"`
	Transform           string          `json:"transform"`
	ColumnIndex         string          `json:"column_index"`
	StartDate           string          `json:"start_date"`
	EndDate             string          `json:"end_date"`
	Data                [][]interface{} `json:"data"`
	Collapse            string          `json:"collapse"`
	Order               string          `json:"order"`
	DatabaseId          int             `json:"database_id"`

	//{"dataset":{"id":4129056,"dataset_code":"EMPL_SEC","database_code":"ADP","name":"Employment by Sector","description":"The ADP National Employment Report (R) is published monthly by the ADP Research Institute (SM) in close collaboration with Moody's Analytics and its experienced team of labor market researchers. The ADP National Employment Report provides a monthly snapshot of U.S. nonfarm private sector employment based on actual transactional payroll data.","refreshed_at":"2015-10-01T07:44:12.207Z","newest_available_date":"2015-09-30","oldest_available_date":"2001-04-30","column_names":["Month","Total Private","Goods Producing","Service Providing"],"frequency":"monthly","type":"Time Series","premium":false,"limit":null,"transform":null,"column_index":null,"start_date":"2001-04-30","end_date":"2015-09-30","data":["2001-04-30",111573.0,24275.29416846,87297.757624888]],"collapse":null,"order":"desc","database_id":589}}
}

type HavenDocumentQuery struct {
	Documents []HavenDocument `json:"documents"`
}

type HavenDocument struct {
	Reference string   `json:"reference"`
	Weight    float32  `json:"weight"`
	Links     []string `json:"links"`
	Index     string   `json:"index"`
	Title     string   `json:"title"`
	//[{"reference":"http://feeds.latimes.com/~r/latimes/news/politics/~3/rHULjutD3Eg/la-pn-rand-paul-mcconnell-20150602-story.html","weight":84.55,"links":["KENTUCK"],"index":"news_eng","title":"Rand Paul and Mitch McConnell alliance of convenience is put to the test"}]
}

func getQuandlData(rw http.ResponseWriter, req *http.Request) {
	var urlVal string
	urlParms, _ := getUrlParameters(req.URL.String())

	for key, _ := range stateList {
		urlVal = fmt.Sprintf(stateUnemploymentRequestTemplate, key, urlParms["from"][0], urlParms["to"][0])
		var tmp QuandlEmploymentData
		err := getJson(urlVal, &tmp)

		if err != nil {
			log.Println("Error Occured")
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("Data Received for state [%s] url [%s] ", key, urlVal))
			log.Println(tmp)
		}

		stateRecs := make([]StateUnemploymentDataRecord, len(tmp.Dataset.Data))
		for i, rec := range tmp.Dataset.Data {
			stateRecs[i] = StateUnemploymentDataRecord{rec[0].(string), rec[1].(float64)}
		}

		stateData[key] = stateRecs
	}

	json.NewEncoder(rw).Encode(stateData)
	fmt.Fprintf(rw, req.URL.Path)
	//	tmp, err := getRawBody(url)
}

func getHavenData(rw http.ResponseWriter, req *http.Request) {
	urlParms, _ := getUrlParameters(req.URL.String())

	doneChannel := make(chan bool)
	doneCount := len(stateList)
	for key, value := range stateList {
		go processHavenData(doneChannel, key, value, urlParms)
	}

	for doneCount > 0 {
		<-doneChannel
		doneCount--
	}
	close(doneChannel)

	json.NewEncoder(rw).Encode(stateSentimentData)

	//	tmp, err := getRawBody(url)

}

func processHavenData(done chan bool, key string, value string, urlParms map[string][]string) {
    urlVal := fmt.Sprintf(havenNewQueryTemplate, value, urlParms["from"][0])
	var tmp HavenDocumentQuery
	err := getJson(urlVal, &tmp)
	//val, err := getRawBody(urlVal)
	if err != nil {
		log.Println("Error Occured")
		log.Println(err)
	} else {
		log.Println(fmt.Sprintf("Data Received for state [%s] url [%s] ", key, urlVal))
		log.Println(tmp)
	}

	stateRecs := make(map[string]StateArticleSentimentRecord)
	for _, rec := range tmp.Documents {
		stateRecs[rec.Title] = StateArticleSentimentRecord{ArticleLink: rec.Reference}
	}

	stateSentimentData[key] = stateRecs

	done <- true
}

func getUrlParameters(urlVal string) (map[string][]string, error) {
	log.Println("Parsing URL: ", urlVal)
	u, err := url.Parse(urlVal)
	if err != nil {
		return nil, err
	}

	return url.ParseQuery(u.RawQuery)
}

func getRawBody(urlVal string) (string, error) {
	r, err := http.Get(urlVal)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	output, _ := ioutil.ReadAll(r.Body)
	return string(output), nil
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&target)
}

func main() {
	log.Println("initializing...")
	initialize()
	http.HandleFunc("/quandl", getQuandlData)
	http.HandleFunc("/haven", getHavenData)

	log.Println("Listening and running")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func initialize() {
	stateData = make(map[string][]StateUnemploymentDataRecord)
	stateSentimentData = make(map[string](map[string]StateArticleSentimentRecord))
}
