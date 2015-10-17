package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func getQuandlData(rw http.ResponseWriter, req *http.Request) {
	url := "https://www.quandl.com/api/v3/datasets/ADP/EMPL_SEC.json?auth_token=7GKGd3eZ4wexWY4Ge1bb&start_date=2015-09-01&end_date=2015-09-30"
	//	tmp, err := getRawBody(url)
	var tmp QuandlEmploymentData
	err := getJson(url, &tmp)

	if err != nil {
		log.Println("Error Occured")
		log.Println(err)
	} else {
		log.Println("Data Received")
		log.Println(tmp)
		//fmt.Fprintf(rw, tmp)
		//rw.Println(tmp)
	}
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
	http.HandleFunc("/quandl", getQuandlData)

	log.Fatal(http.ListenAndServe(":8082", nil))
}

//{"dataset":{"id":4129056,"dataset_code":"EMPL_SEC","database_code":"ADP","name":"Employment by Sector","description":"The ADP National Employment Report (R) is published monthly by the ADP Research Institute (SM) in close collaboration with Moody's Analytics and its experienced team of labor market researchers. The ADP National Employment Report provides a monthly snapshot of U.S. nonfarm private sector employment based on actual transactional payroll data.","refreshed_at":"2015-10-01T07:44:12.207Z","newest_available_date":"2015-09-30","oldest_available_date":"2001-04-30","column_names":["Month","Total Private","Goods Producing","Service Providing"],"frequency":"monthly","type":"Time Series","premium":false,"limit":null,"transform":null,"column_index":null,"start_date":"2001-04-30","end_date":"2015-09-30","data":["2001-04-30",111573.0,24275.29416846,87297.757624888]],"collapse":null,"order":"desc","database_id":589}}
