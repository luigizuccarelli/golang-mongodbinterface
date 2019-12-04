package main

import (
	"crypto/tls"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"time"
)

type Clients interface {
	DBGetAffiliates() ([]Affiliate, error)
	DBGetPublications(string, []byte) ([]Publication, error)
	DBGetStocks(string, bool) ([]Stock, error)
	DBUpdateStock(body []byte) ([]Stock, error)
	DBIndex() error
	DBSetup(body []byte) error
	DBMigrate(body []byte) error
	DBUpdateAffiliateSpecific(body []byte) error
	DBUpdateStockCurrentPrice() error
	DBUpdateWatchlist(body []byte) (Watchlist, error)
	DBGetWatchlist(id string) (Watchlist, error)
	DBGetStocksCount(id string) (int, error)
	DBGetStocksPaginated(id string, skip int, limit int) ([]Stock, error)
	GetPriceStatus() (string, error)
	SendAlert([]byte) error
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Close() error
}

type Connectors struct {
	session *mgo.Session
	http    *http.Client
	redis   *redis.Client
	name    string
}

func NewClientConnectors() Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	// mongodb connection
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("MONGODB_HOST") + ":" + os.Getenv("MONGODB_PORT")},
		Timeout:  40 * time.Second,
		Database: os.Getenv("MONGODB_DATABASE"),
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	// database setup and init
	s, e := mgo.DialWithInfo(mongoDBDialInfo)
	if e != nil {
		logger.Error(fmt.Sprintf("Mongodb init %v ", e.Error()))
		return nil
	}
	s.SetMode(mgo.Monotonic, true)
	logger.Trace(fmt.Sprintf("Mongodb connection successful %v ", s))

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
	})

	return &Connectors{http: httpClient, redis: redisClient, session: s, name: "RealConnectors"}
}

func (r *Connectors) Get(key string) (string, error) {
	val, err := r.redis.Get(key).Result()
	return val, err
}

func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	val, err := r.redis.Set(key, value, expr).Result()
	return val, err
}

func (r *Connectors) Close() error {
	r.redis.Close()
	return nil
}
