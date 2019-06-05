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
	Version   string `json:"version"`
	Level     string `json:"level"`
	Basedir   string `json:"base_dir"`
	Provider  string `json:"provider"`
	Providers []Provider
	Port      string `json:"port"`
	Cache     string `json:"cache"`
	Url       string `json:"url"`
	QuoteUrl  string `json:"quoteurl"`
	Token     string `json:"token"`
	MongoDB   Mongodb
	RedisDB   Redis
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

type Provider struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type Subtrades struct {
	SubstradeSetting `json:"subtradesetting"`
}

type SubstradeSetting struct {
	Stop float64 `json:"smartstop"`
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
	CurrencySign   string        `json:"currencysign"`
}

type Watchlist struct {
	UID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	CustomerId int           `json:"customerid"`
	Stocks     []string      `json:"stocks"`
}

type TradeSmithSchema struct {
	Symbol         string          `json:"symbol"`
	Buy            float64         `json:"openprice"`
	SubTrades      []Subtrades     `json:"subtrades"`
	Status         int             `json:"tradestatus"`
	Recommendation PositionSetting `json:"positionsetting"`
	CurrencySign   string          `json:"currencysygn"`
}

type PositionSetting struct {
	Info  string `json:"text1"`
	Other string `json:"text2"`
}

// alpha advantage GLOBAL_QUOTE struct
type Alphavantage struct {
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

type IEXCloud struct {
	Symbol                string      `json:"symbol"`
	CompanyName           string      `json:"companyName"`
	CalculationPrice      string      `json:"calculationPrice"`
	Open                  interface{} `json:"open"`
	OpenTime              interface{} `json:"openTime"`
	Close                 float64     `json:"close"`
	CloseTime             int64       `json:"closeTime"`
	High                  interface{} `json:"high"`
	Low                   interface{} `json:"low"`
	LatestPrice           float64     `json:"latestPrice"`
	LatestSource          string      `json:"latestSource"`
	LatestTime            string      `json:"latestTime"`
	LatestUpdate          int64       `json:"latestUpdate"`
	LatestVolume          int         `json:"latestVolume"`
	IexRealtimePrice      interface{} `json:"iexRealtimePrice"`
	IexRealtimeSize       interface{} `json:"iexRealtimeSize"`
	IexLastUpdated        interface{} `json:"iexLastUpdated"`
	DelayedPrice          interface{} `json:"delayedPrice"`
	DelayedPriceTime      interface{} `json:"delayedPriceTime"`
	ExtendedPrice         float64     `json:"extendedPrice"`
	ExtendedChange        float64     `json:"extendedChange"`
	ExtendedChangePercent float64     `json:"extendedChangePercent"`
	ExtendedPriceTime     int64       `json:"extendedPriceTime"`
	PreviousClose         float64     `json:"previousClose"`
	Change                int         `json:"change"`
	ChangePercent         int         `json:"changePercent"`
	IexMarketPercent      interface{} `json:"iexMarketPercent"`
	IexVolume             interface{} `json:"iexVolume"`
	AvgTotalVolume        int         `json:"avgTotalVolume"`
	IexBidPrice           interface{} `json:"iexBidPrice"`
	IexBidSize            interface{} `json:"iexBidSize"`
	IexAskPrice           interface{} `json:"iexAskPrice"`
	IexAskSize            interface{} `json:"iexAskSize"`
	MarketCap             int64       `json:"marketCap"`
	PeRatio               float64     `json:"peRatio"`
	Week52High            float64     `json:"week52High"`
	Week52Low             int         `json:"week52Low"`
	YtdChange             float64     `json:"ytdChange"`
}

// ShcemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
	ID           bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	LastUpdate   int64         `json:"lastupdate,omitempty"`
	MetaInfo     string        `json:"metainfo,omitempty"`
	Count        int64         `json:"count,omitempty"`
	TotalPages   int64         `json:"totalpages,omitempty"`
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
