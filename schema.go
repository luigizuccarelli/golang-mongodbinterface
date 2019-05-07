package main

import (
	"gopkg.in/mgo.v2/bson"
)

/**
 * The json config will always be in the form of parent.key:value pair
 * the reasoning here is that it is easy to maintain and use
 * and also if required can be migtrated to a key value store such as redis
 *
 * Don't dig it - then feel welcome to change it to your hearts content - knock yourself out
 *
 **/

// Config structure - define the json format for our microservice config
type Config struct {
	Version  string `json:"version"`
	Level    string `json:"level"`
	Basedir  string `json:"base_dir"`
	Port     string `json:"port"`
	Cache    string `json:"cache"`
	Url      string `json:"url"`
	QuoteUrl string `json:"quoteurl"`
	Token    string `json:"token"`
	MongoDB  Mongodb
	RedisDB  Redis
}

// Mongodb structure - the base config to connect to mongodb
type Mongodb struct {
	Host           string `json:"host"`
	Port           string `json:"port"`
	DatabaseName   string `json:"name"`
	User           string `json:"user"`
	Password       string `json:"pwd"`
	AdminUser      string `json:"adminuser"`
	AdminPasssword string `json:"adminpwd"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Affiliate struct {
	UID   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Id    int           `json:"id"`
	Name  string        `json:"name"`
	Token string        `json:"token"`
}

type Publication struct {
	UID         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	AffiliateId int           `json:"affiliateid"`
}

type Stock struct {
	UID            bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	PublicationId  int           `json:"publicationid"`
	AffiliateId    int           `json:"affiliateid"`
	RefId          int           `json:"id"`
	Symbol         string        `json:"symbol"`
	Name           string        `json:"name"`
	Buy            float64       `json:"buy"`
	Stop           float64       `json:"stop"`
	Last           float64       `json:"last"`
	Change         float64       `json:"change"`
	Recommendation string        `json:"recommendation"`
	Status         int           `json:"status"`
}

type Watchlist struct {
	UID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	CustomerID string        `json:"customerid"`
	Stocks     []string      `json:"stocks"`
}

type TradeSmithSchema struct {
	Symbol         string          `json:"symbol"`
	Buy            float64         `json:"openprice"`
	Stop           float64         `json:"stoprice"`
	Status         int             `json:"tradestatus"`
	Recommendation PositionSetting `json:"positionsetting"`
}

type PositionSetting struct {
	Info  string `json:"text1"`
	Other string `json:"text2"`
}

// alpha advantage GLOBAL_QUOTE struct
type AlphaAdvantage struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Open             string `json:"02. open"`
		High             string `json:"03. high"`
		Low              string `json:"04. low"`
		Price            string `json:"05. price"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
		PreviousClose    string `json:"08. previous close"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
	} `json:"Global Quote"`
}

// ShcemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
	ID           bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	LastUpdate   int64         `json:"lastupdate,omitempty"`
	MetaInfo     string        `json:"metainfo,omitempty"`
	Affiliates   []Affiliate   `json:"affiliates"`
	Publications []Publication `json:"publications"`
	Stocks       []Stock       `json:"stocks"`
	WatchList    Watchlist     `json:"watchlist"`
}

// Response schema
type Response struct {
	StatusCode string          `json:"statuscode"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Payload    SchemaInterface `json:"payload"`
}
