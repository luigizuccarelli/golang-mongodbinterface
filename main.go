package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	logger     simple.Logger
	connectors Clients
)

func startHttpServer() *http.Server {

	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()
	r.HandleFunc("/api/v2/sys/info/isalive", IsAlive).Methods("GET")
	r.HandleFunc("/api/v1/migrate", MiddlewareMigrateData).Methods("POST")
	r.HandleFunc("/api/v1/setup", MiddlewareDBSetup).Methods("POST")
	r.HandleFunc("/api/v1/index", MiddlewareDBIndex).Methods("POST")
	r.HandleFunc("/api/v1/specific", MiddlewareUpdateSpecific).Methods("POST")
	r.HandleFunc("/api/v1/affiliates", MiddlewareDBGetAllAffiliates).Methods("GET")
	r.HandleFunc("/api/v1/publications/{affiliateid}", MiddlewareDBGetAllPublicationsByAffiliate).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/v1/stocks/{publicationid}", MiddlewareDBGetStocksByPublication).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/affiliate/{affiliateid}", MiddlewareDBGetAllStocksByAffiliatePaginated).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/{bsonid}", MiddlewareDBUpdateStock).Methods("OPTIONS", "POST", "PUT")
	r.HandleFunc("/api/v1/stocks/affiliate/{affiliateid}/count", MiddlewareDBGetAllStocksCount).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/watchlist/{customerid}", MiddlewareDBGetWatchlist).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/watchlist/{customerid}", MiddlewareDBUpdateWatchlist).Methods("OPTIONS", "PUT", "POST")
	r.HandleFunc("/api/v1/prices", MiddlewareDBUpdateStockCurrentPrice).Methods("POST")
	r.HandleFunc("/api/v1/pricestatus", MiddlewarePriceStatus).Methods("GET")

	sh := http.StripPrefix("/api/v2/api-docs/", http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/v2/api-docs/").Handler(sh)
	http.Handle("/", r)

	connectors = NewClientConnectors()

	// lightweight thread of execution
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {
	err := ValidateEnvars()
	if err != nil {
		os.Exit(-1)
	}
	logger.Level = os.Getenv("LOG_LEVEL")
	srv := startHttpServer()
	logger.Info("Starting server on port " + os.Getenv("SERVER_PORT"))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}
