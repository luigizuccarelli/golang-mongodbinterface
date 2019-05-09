package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger     simple.Logger
	config     Config
	connectors Clients
)

func startHttpServer(cfg Config) *http.Server {

	config = cfg

	logger.Debug(fmt.Sprintf("Config in startServer %v ", cfg))
	srv := &http.Server{Addr: ":" + cfg.Port}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sys/info/isalive", IsAlive).Methods("GET")
	r.HandleFunc("/api/v1/migrate", MiddlewareMigrateData).Methods("POST")
	r.HandleFunc("/api/v1/setup", MiddlewareDBSetup).Methods("POST")
	r.HandleFunc("/api/v1/index", MiddlewareDBIndex).Methods("POST")
	r.HandleFunc("/api/v1/specific", MiddlewareUpdateSpecific).Methods("POST")
	r.HandleFunc("/api/v1/affiliates", MiddlewareDBGetAllAffiliates).Methods("GET")
	r.HandleFunc("/api/v1/publications/{affiliateid}", MiddlewareDBGetAllPublicationsByAffiliate).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/{publicationid}", MiddlewareDBGetStocksByPublication).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/affiliate/{affiliateid}", MiddlewareDBGetAllStocksByAffiliate).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/{bsonid}", MiddlewareDBUpdateStock).Methods("OPTIONS", "POST", "PUT")
	r.HandleFunc("/api/v1/watchlist/{customerid}", MiddlewareDBGetWatchlist).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/watchlist/{bsonid}", MiddlewareDBUpdateWatchlist).Methods("OPTIONS", "PUT", "POST")
	r.HandleFunc("/api/v1/prices", MiddlewareDBUpdateStockCurrentPrice).Methods("POST")
	http.Handle("/", r)

	connectors = NewClientConnectors(cfg)

	// lightweight thread of execution
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {
	// read the config
	config, _ := Init("config.json")
	logger.Level = config.Level
	srv := startHttpServer(config)
	logger.Info("Starting server on port " + config.Port)
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
