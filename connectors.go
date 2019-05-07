package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type Clients interface {
	DBGetAffiliates() ([]Affiliate, error)
	DBGetPublications(string) ([]Publication, error)
	DBGetStocks(string, bool) ([]Stock, error)
	DBUpdateStock(body []byte) ([]Stock, error)
	DBIndex() error
	DBSetup(body []byte) error
	DBMigrate(body []byte) error
	DBUpdateAffiliateSpecific(body []byte) error
	DBUpdateStockCurrentPrice() error
	DBUpdateWatchlist(body []byte) (Watchlist, error)
	DBGetWatchlist(id string) (Watchlist, error)
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

func NewClientConnectors(cfg Config) Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	// mongodb connection
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{cfg.MongoDB.Host + ":" + cfg.MongoDB.Port},
		Timeout:  40 * time.Second,
		Database: cfg.MongoDB.DatabaseName,
		Username: cfg.MongoDB.User,
		Password: cfg.MongoDB.Password,
	}

	// database setup and init
	s, e := mgo.DialWithInfo(mongoDBDialInfo)
	if e != nil {
		logger.Error(fmt.Sprintf("Mongodb init %v ", e.Error()))
		return nil
	}
	s.SetMode(mgo.Monotonic, true)
	logger.Trace(fmt.Sprintf("Mongodb connection successful %v ", s))

	return &Connectors{http: httpClient, session: s, name: "RealConnectors"}
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
