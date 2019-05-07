package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// MiddlewareDBSetup a http response and request wrapper for portfolio's that are associated to affiliate
// It takes a both response and request objects and returns void
func MiddlewareDBSetup(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareDBSetup) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = connectors.DBSetup(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data migrate (MiddlewareDBSetup) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database setup"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBSetup) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBIndex a http response and request wrapper for portfolio's that are associated to affiliate
// It takes a both response and request objects and returns void
func MiddlewareDBIndex(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	err := connectors.DBIndex()
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBIndex) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database index"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBIndex) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetAllAffiliates a http response and request wrapper for affiliate data
// It takes a both response and request objects and returns void
func MiddlewareDBGetAllAffiliates(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	affiliates, err := connectors.DBGetAffiliates()
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBGetAllAffiliates) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to affiliates", Affiliates: affiliates}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllAffiliates) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetAllPublicationsByAffiliate a http response and request wrapper for affiliate data
// It takes a both response and request objects and returns void
func MiddlewareDBGetAllPublicationsByAffiliate(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	vars := mux.Vars(r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	publications, err := connectors.DBGetPublications(vars["affiliateid"])
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBGetAllPublicationsByAffiliate) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to publications for affiliate " + vars["affiliateid"], Publications: publications}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllPublicationsByAffiliate) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetStocksByPublication a http response and request wrapper for stock data
// It takes a both response and request objects and returns void
func MiddlewareDBGetStocksByPublication(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	vars := mux.Vars(r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	stocks, err := connectors.DBGetStocks(vars["publicationid"], false)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBGetAllStocks) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to stocks " + vars["publicationid"], Stocks: stocks}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllStocks) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetAllStocksByAffiliate a http response and request wrapper for stock data
// It takes a both response and request objects and returns void
func MiddlewareDBGetAllStocksByAffiliate(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	vars := mux.Vars(r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	stocks, err := connectors.DBGetStocks(vars["affiliateid"], true)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBGetAllStocksByAffiliate) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to stocks " + vars["affiliateid"], Stocks: stocks}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllStocksByAffiliate) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareMigrateData a http response and request wrapper for portfolio's that are associated to affiliate
// It takes a both response and request objects and returns void
func MiddlewareMigrateData(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareMigrateData) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = connectors.DBMigrate(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data migrate (MiddlewareMigrateData) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Affiliate data (publication and stocks)"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareMigrateData) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareUpdateSpecific a http response and request wrapper for stocks's that are associated to affiliate
// It takes a both response and request objects and returns void
func MiddlewareUpdateSpecific(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareUpdateSpecific) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = connectors.DBUpdateAffiliateSpecific(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data specific update (MiddlewareUpdateSpecific) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Update of affiliate specific data for all associated stocks"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareUpdateSpecific) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBUpdateStockCurrentPrice a http response and request wrapper to update the stock price and percentage change
// It takes a both response and request objects and returns void
func MiddlewareDBUpdateStockCurrentPrice(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	e := connectors.DBUpdateStockCurrentPrice()
	if e != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data update (MiddlewareDBUpdateStockCurrentPrice) " + e.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Stock update from API"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBUpdateStockCurrentPrice) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBUpdateStock a http response and request wrapper to update the stock
// It takes a both response and request objects and returns void
func MiddlewareDBUpdateStock(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareDBUpdateStock) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	st, e := connectors.DBUpdateStock(body)
	if e != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data update (MiddlewareDBUpdateStock) " + e.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Stock update", Stocks: st}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBUpdateStock) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetWatchlist a http response and request wrapper for customer watchlist data
// It takes a both response and request objects and returns void
func MiddlewareDBGetWatchlist(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	vars := mux.Vars(r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	wl, err := connectors.DBGetWatchlist(vars["customerid"])
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Indexing (MiddlewareDBGetWatchlist) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call for watchlist " + vars["customerid"], WatchList: wl}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetWatchlist) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBUpdateWatchlist a http response and request wrapper to update the stock
// It takes a both response and request objects and returns void
func MiddlewareDBUpdateWatchlist(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	addHeaders(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareDBUpdateWatchlist) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	wl, e := connectors.DBUpdateWatchlist(body)
	if e != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Data update (MiddlewareDBUpdateWatchlist) " + e.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Watchlist update", WatchList: wl}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBUpdateWatchlist) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	logger.Trace(fmt.Sprintf("used to mask cc %v", r))
	logger.Trace(fmt.Sprintf("config data  %v", config))
	fmt.Fprintf(w, "{\"isalive\": true , \"version\": \""+config.Version+"\"}")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	var request []string
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	logger.Trace(fmt.Sprintf("Headers : %s", request))

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}
