package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
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
		response = handleError(w, "Could not read body data (MiddlewareDBSetup) "+err.Error(), payload)
	}

	err = connectors.DBSetup(body)
	if err != nil {
		response = handleError(w, "Data migrate (MiddlewareDBSetup) "+err.Error(), payload)
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
		response = handleError(w, "Indexing (MiddlewareDBIndex) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	affiliates, err := connectors.DBGetAffiliates()
	if err != nil {
		response = handleError(w, "Indexing (MiddlewareDBGetAllAffiliates) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	publications, err := connectors.DBGetPublications(vars[AFFILIATEID])
	if err != nil {
		response = handleError(w, "Indexing (MiddlewareDBGetAllPublicationsByAffiliate) "+err.Error(), payload)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to publications for affiliate " + vars[AFFILIATEID], Publications: publications}
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

	addHeaders(w, r)
	handleOptions(w, r)

	stocks, err := connectors.DBGetStocks(vars["publicationid"], false)
	if err != nil {
		response = handleError(w, "Indexing (MiddlewareDBGetAllStocks) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	stocks, err := connectors.DBGetStocks(vars[AFFILIATEID], true)
	if err != nil {
		response = handleError(w, "MW call (MiddlewareDBGetAllStocksByAffiliate) "+err.Error(), payload)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Database call to stocks " + vars[AFFILIATEID], Stocks: stocks}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllStocksByAffiliate) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetAllStocksCount a http response and request wrapper for stock data
// It takes a both response and request objects and returns void
func MiddlewareDBGetAllStocksCount(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	vars := mux.Vars(r)

	addHeaders(w, r)
	handleOptions(w, r)

	count, err := connectors.DBGetStocksCount(vars[AFFILIATEID])
	if err != nil {
		response = handleError(w, "MW call (MiddlewareDBGetAllStocksCount) "+err.Error(), payload)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Stocks count for affiliateid " + vars[AFFILIATEID] + ":" + strconv.Itoa(count)}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBGetAllStocksCount) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBGetAllStocksByAffiliatePaginated a http response and request wrapper for stock data
// It takes a both response and request objects and returns void
func MiddlewareDBGetAllStocksByAffiliatePaginated(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	var totalPages int64 = 0
	var limit, skip int
	vars := mux.Vars(r)

	addHeaders(w, r)
	handleOptions(w, r)

	offset := r.URL.Query().Get("offset")
	page := r.URL.Query().Get("perpage")

	count, err := connectors.DBGetStocksCount(vars[AFFILIATEID])
	if page != "" && offset != "" {
		total, _ := strconv.Atoi(page)
		totalPages = int64(count / total)
		skip, _ = strconv.Atoi(offset)
		limit, _ = strconv.Atoi(page)
	} else {
		totalPages = (int64(count) / 10)
		skip = 0
		limit = count
	}
	stocks, err := connectors.DBGetStocksPaginated(vars[AFFILIATEID], skip, limit)
	if err != nil {
		response = handleError(w, "Paginated stocks (MiddlewareDBGetAllStocksByAffiliate) "+err.Error(), payload)
	} else {
		payload = SchemaInterface{
			LastUpdate: time.Now().Unix(),
			MetaInfo:   "MiddlewareDBGetAllStocksByAffiliate for affiliateId " + vars[AFFILIATEID],
			Count:      int64(count),
			TotalPages: totalPages,
			Stocks:     stocks,
		}
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

	addHeaders(w, r)
	handleOptions(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = handleError(w, "Could not read body data (MiddlewareMigrateData) "+err.Error(), payload)
	}

	err = connectors.DBMigrate(body)
	if err != nil {
		response = handleError(w, "Data migrate (MiddlewareMigrateData) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = handleError(w, "Could not read body data (MiddlewareUpdateSpecific) "+err.Error(), payload)
	}

	err = connectors.DBUpdateAffiliateSpecific(body)
	if err != nil {
		response = handleError(w, "Data specific update (MiddlewareUpdateSpecific) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	e := connectors.DBUpdateStockCurrentPrice()
	if e != nil {
		response = handleError(w, "Data update (MiddlewareDBUpdateStockCurrentPrice) "+e.Error(), payload)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Stock update from API Status : PENDING"}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBUpdateStockCurrentPrice) status : PENDING", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareDBUpdateStock a http response and request wrapper to update the stock
// It takes a both response and request objects and returns void
func MiddlewareDBUpdateStock(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	addHeaders(w, r)
	handleOptions(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = handleError(w, "Could not read body data (MiddlewareDBUpdateStock) "+err.Error(), payload)
	}

	st, e := connectors.DBUpdateStock(body)
	if e != nil {
		response = handleError(w, "Data update (MiddlewareDBUpdateStock) "+e.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	wl, err := connectors.DBGetWatchlist(vars["customerid"])
	if err != nil {
		response = handleError(w, "Querying (MiddlewareDBGetWatchlist) "+err.Error(), payload)
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

	addHeaders(w, r)
	handleOptions(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = handleError(w, "Could not read body data (MiddlewareDBUpdateWatchlist)"+err.Error(), payload)
	}

	wl, e := connectors.DBUpdateWatchlist(body)
	if e != nil {
		response = handleError(w, "Data update (MiddlewareDBUpdateWatchlist) "+e.Error(), payload)
	} else {
		payload = SchemaInterface{LastUpdate: time.Now().Unix(), MetaInfo: "Watchlist update", WatchList: wl}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW call (MiddlewareDBUpdateWatchlist) successfull", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func MiddlewarePriceStatus(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	val, _ := connectors.GetPriceStatus()
	logger.Info(fmt.Sprintf("Price update status : %s", val))
	fmt.Fprintf(w, "{\"status\": \""+val+"\"}\n")
}

// IsAlive - liveliness and readiness probe check
func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	logger.Trace(fmt.Sprintf("used to mask cc %v", r))
	logger.Trace(fmt.Sprintf("config data  %v", config))
	fmt.Fprintf(w, "{\"isalive\": true , \"version\": \""+config.Version+"\"}\n")
}

// simple options handler
func handleOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	}
	return
}

// simple error handler
func handleError(w http.ResponseWriter, msg string, s SchemaInterface) Response {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{StatusCode: "500", Status: "ERROR", Message: msg, Payload: s}
	return r
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
