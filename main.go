package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var stateUnemploymentRequestTemplate string = "https://www.quandl.com/api/v3/datasets/FRED/%sUR.json?auth_token=7GKGd3eZ4wexWY4Ge1bb&start_date=%s&end_date=%s"

type StateUnemploymentDataRecord struct {
	Date  string
	Value float64
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
}

var stateData map[string]([]StateUnemploymentDataRecord)

var stateList = []string{
"AL","AK","AZ","AR",
"CA","CO","CT","DE",
"FL","GA","HI","ID",
"IL","IN","IA","KS",
"KY","LA","ME","MD",
"MA","MI","MN","MS",
"MO","MT","NE","NV",
"NH","NJ","NM","NY",
"NC","ND","OH","OK",
"OR","PA","RI","SC",
"SD","TN","TX","UT",
"VT","VA","WA","WV",
"WI","WY","DC"}

func getQuandlData(rw http.ResponseWriter, req *http.Request) {
	var url string
	for _, state := range stateList {
		url = fmt.Sprintf(stateUnemploymentRequestTemplate, state, "2001-01-01", "2015-12-30")
		var tmp QuandlEmploymentData
		err := getJson(url, &tmp)

		if err != nil {
			log.Println("Error Occured")
			log.Println(err)
		} else {
            log.Println(fmt.Sprintf("Data Received for state [%s] url [%s] ", state, url))
			log.Println(tmp)
		}
        
        
        stateRecs := make([]StateUnemploymentDataRecord, len(tmp.Dataset.Data))
        for i, rec := range tmp.Dataset.Data {
            stateRecs[i] = StateUnemploymentDataRecord {rec[0].(string), rec[1].(float64)}    
        }
        
        stateData[state] = stateRecs
	}
    
    json.NewEncoder(rw).Encode(stateData)
	//	tmp, err := getRawBody(url)

}

func getRawBody(url string) (string, error) {
	r, err := http.Get(url)
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
    initialize()
	http.HandleFunc("/quandl", getQuandlData)

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func initialize() {
	stateData = make(map[string][]StateUnemploymentDataRecord)
}

//{"dataset":{"id":4129056,"dataset_code":"EMPL_SEC","database_code":"ADP","name":"Employment by Sector","description":"The ADP National Employment Report (R) is published monthly by the ADP Research Institute (SM) in close collaboration with Moody's Analytics and its experienced team of labor market researchers. The ADP National Employment Report provides a monthly snapshot of U.S. nonfarm private sector employment based on actual transactional payroll data.","refreshed_at":"2015-10-01T07:44:12.207Z","newest_available_date":"2015-09-30","oldest_available_date":"2001-04-30","column_names":["Month","Total Private","Goods Producing","Service Providing"],"frequency":"monthly","type":"Time Series","premium":false,"limit":null,"transform":null,"column_index":null,"start_date":"2001-04-30","end_date":"2015-09-30","data":["2001-04-30",111573.0,24275.29416846,87297.757624888]],"collapse":null,"order":"desc","database_id":589}}
