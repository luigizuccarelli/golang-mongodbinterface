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

func checkEnvar(item string) {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	if os.Getenv(name) == "" {
		if required {
			logger.Error(fmt.Sprintf("%s envar is mandatory please set it", name))
			os.Exit(-1)
		} else {
			logger.Error(fmt.Sprintf("%s envar is empty please set it", name))
		}
	}
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
func ValidateEnvars() {
	items := []string{
		"LOG_LEVEL,false",
		"SERVER_PORT,true",
		"REDIS_HOST,true",
		"REDIS_PORT,true",
		"REDIS_PASSWORD,true",
		"MONGODB_HOST,true",
		"MONGODB_DATABASE,true",
		"MONGODB_USER,true",
		"MONGODB_PASSWORD,true",
		"VERSION,true",
		"URL,true",
		"PROVIDER_NAME,true",
		"PROVIDER_URL,true",
		"PROVIDER_TOKEN,true",
		"ANALYTICS_URL,true",
	}
	for x, _ := range items {
		checkEnvar(items[x])
	}
}
