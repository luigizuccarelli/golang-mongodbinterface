package main

import (
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger     simple.Logger
	connectors Clients
)

func startHttpServer() *http.Server {

	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sys/info/isalive", IsAlive).Methods("GET")
	r.HandleFunc("/api/v1/migrate", MiddlewareMigrateData).Methods("POST")
	r.HandleFunc("/api/v1/setup", MiddlewareDBSetup).Methods("POST")
	r.HandleFunc("/api/v1/index", MiddlewareDBIndex).Methods("POST")
	r.HandleFunc("/api/v1/specific", MiddlewareUpdateSpecific).Methods("POST")
	r.HandleFunc("/api/v1/affiliates", MiddlewareDBGetAllAffiliates).Methods("GET")
	r.HandleFunc("/api/v1/publications/{affiliateid}", MiddlewareDBGetAllPublicationsByAffiliate).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/{publicationid}", MiddlewareDBGetStocksByPublication).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/affiliate/{affiliateid}", MiddlewareDBGetAllStocksByAffiliatePaginated).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/stocks/{bsonid}", MiddlewareDBUpdateStock).Methods("OPTIONS", "POST", "PUT")
	r.HandleFunc("/api/v1/stocks/affiliate/{affiliateid}/count", MiddlewareDBGetAllStocksCount).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/watchlist/{customerid}", MiddlewareDBGetWatchlist).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v1/watchlist/{customerid}", MiddlewareDBUpdateWatchlist).Methods("OPTIONS", "PUT", "POST")
	r.HandleFunc("/api/v1/prices", MiddlewareDBUpdateStockCurrentPrice).Methods("POST")
	r.HandleFunc("/api/v1/pricestatus", MiddlewarePriceStatus).Methods("GET")

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/swaggerui/").Handler(sh)
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
	ValidateEnvars()
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

func ValidateEnvars() {
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "info")
	}
	if os.Getenv("SERVER_PORT") == "" {
		os.Setenv("SERVER_PORT", "9000")
	}
	if os.Getenv("MONGODB_HOST") == "" {
		logger.Error("MONGODB_HOST envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("MONGODB_DATABASE") == "" {
		logger.Error("MONGODB_DATABASE envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("MONGODB_USER") == "" {
		logger.Error("MONGODB_USER envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("MONGODB_PASSWORD") == "" {
		logger.Error("MONGODB_PASSWORD envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("REDIS_HOST") == "" {
		logger.Error("MONGODB_HOST envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("REDIS_PORT") == "" {
		logger.Error("REDIS_PORT envar is mandatory")
		os.Exit(-1)
	}
	if os.Getenv("REDIS_PASSWORD") == "" {
		logger.Error("REDIS_PASSWORD envar is mandatory")
		os.Exit(-1)
	}
}
