package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllMiddleware(t *testing.T) {
	var req *http.Request
	var response Response

	// create anonymous struct
	tests := []struct {
		Name     string
		Method   string
		Url      string
		Payload  string
		Handler  string
		FileName string
		want     int
		errorMsg string
	}{
		{
			"Test [Isalive] should pass",
			"GET", "api/v1/sys/info/isalive",
			"",
			"IsAlive",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBSetup] should pass",
			"POST",
			"api/v1/setup",
			"[{\"id\": 1, \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}]",
			"MiddlewareDBSetup",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBSetup] should pass",
			"POST",
			"api/v1/setup",
			"{\"user\": \"\",\"password\":\"\"}",
			"MiddlewareDBSetup",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBIndex] should pass",
			"POST",
			"api/v1/index",
			"{\"username\": \"\",\"password\":\"\"}",
			"MiddlewareDBIndex",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBGetAllAffiliates] should pass",
			"POST",
			"api/v1/affiliates/1",
			"{\"username\": \"\",\"password\":\"\"}",
			"MiddlewareDBGetAllAffiliates",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBGetAllPublicationsByAffiliate] should pass",
			"POST",
			"api/v1/affiliates/1",
			"{\"affiliateid\": \"\",\"2\":\"\"}",
			"MiddlewareDBGetAllPublicationsByAffiliate",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBGetAllStocksByAffiliate] should pass",
			"POST",
			"api/v1/stocks/1/all",
			"{\"username\": \"\",\"password\":\"\"}",
			"MiddlewareDBGetAllStocksByAffiliate",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBGetStocksByPublication] should pass",
			"POST",
			"api/v1/stocks/1",
			"{\"username\": \"\",\"password\":\"\"}",
			"MiddlewareDBGetStocksByPublication",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},

		{
			"[MiddlewareMigrateData] should pass",
			"POST",
			"api/v1/migrate",
			"{\"id\": 1, \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}",
			"MiddlewareMigrateData",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareMigrateData] should fail",
			"POST",
			"api/v1/migrate",
			"[{\"id\": 1, \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}]",
			"MiddlewareMigrateData",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareUpdateSpecific] should pass",
			"POST",
			"api/v1/migrate",
			"{\"id\": 1, \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}",
			"MiddlewareUpdateSpecific",
			"tests/tss.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBUpdateStockCurrentPrice] should pass",
			"POST",
			"api/v1/stocks",
			"",
			"MiddlewareDBUpdateStockCurrentPrice",
			"tests/tss.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBUpdateStock] should pass",
			"POST",
			"api/v1/stocks",
			"{\"_id\": \"5cc042307ccc69ada893144c\", \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}",
			"MiddlewareDBUpdateStock",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBUpdateStock] should pass",
			"POST",
			"api/v1/stocks",
			"{\"_id\": \"na\", \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}",
			"MiddlewareDBUpdateStock",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBGetWatchlist] should pass",
			"POST",
			"api/v1/stocks",
			"{\"_id\": \"na\", \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}",
			"MiddlewareDBGetWatchlist",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewareDBUpdateWatchlist] should pass",
			"POST",
			"api/v1/watchlist",
			"{\"_id\": \"5cc042307ccc69ada893144c\", \"customerid\":13223,\"stocks\": [\"ABC\",\"DEF\"]}",
			"MiddlewareDBUpdateWatchlist",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[MiddlewarePriceStatus] should pass",
			"GET",
			"api/v1/pricestatus",
			"",
			"MiddlewarePriceStatus",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler %s returned - got (%v) wanted (%v)",
		},
	}

	for _, tt := range tests {
		logger.Info(fmt.Sprintf("Executing test : %s \n", tt.Name))
		if tt.Payload == "" {
			req, _ = http.NewRequest(tt.Method, tt.Url, nil)
		} else {
			req, _ = http.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Payload)))
		}

		connectors = NewTestClients(tt.FileName, tt.want)

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		switch tt.Handler {
		case "IsAlive":
			handler := http.HandlerFunc(IsAlive)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBSetup":
			handler := http.HandlerFunc(MiddlewareDBSetup)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBIndex":
			handler := http.HandlerFunc(MiddlewareDBIndex)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBGetAllAffiliates":
			handler := http.HandlerFunc(MiddlewareDBGetAllAffiliates)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBGetAllPublicationsByAffiliate":
			handler := http.HandlerFunc(MiddlewareDBGetAllPublicationsByAffiliate)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBGetAllStocksByAffiliate":
			handler := http.HandlerFunc(MiddlewareDBGetAllStocksByAffiliate)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBGetStocksByPublication":
			handler := http.HandlerFunc(MiddlewareDBGetStocksByPublication)
			handler.ServeHTTP(rr, req)
		case "MiddlewareMigrateData":
			handler := http.HandlerFunc(MiddlewareMigrateData)
			handler.ServeHTTP(rr, req)
		case "MiddlewareUpdateSpecific":
			handler := http.HandlerFunc(MiddlewareUpdateSpecific)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBUpdateStockCurrentPrice":
			handler := http.HandlerFunc(MiddlewareDBUpdateStockCurrentPrice)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBUpdateStock":
			handler := http.HandlerFunc(MiddlewareDBUpdateStock)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBGetWatchlist":
			handler := http.HandlerFunc(MiddlewareDBGetWatchlist)
			handler.ServeHTTP(rr, req)
		case "MiddlewareDBUpdateWatchlist":
			handler := http.HandlerFunc(MiddlewareDBUpdateWatchlist)
			handler.ServeHTTP(rr, req)
		case "MiddlewarePriceStatus":
			handler := http.HandlerFunc(MiddlewarePriceStatus)
			handler.ServeHTTP(rr, req)
		}
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf(fmt.Sprintf(tt.errorMsg, tt.Handler, "nil", "error"))
		}
		// ignore errors here
		json.Unmarshal(body, &response)
		_ = response.Payload.MetaInfo
		if rr.Code != tt.want {
			t.Fatalf(fmt.Sprintf(tt.errorMsg, tt.Handler, "nil", "error"))
		}
	}
}
