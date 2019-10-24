package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/microlib/simple"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	// create a key value map (to fake redis)
	store      map[string]string
	logger     simple.Logger
	config     Config
	connectors Clients
	counter    int = 0
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
	DBGetWatchlist(string) (Watchlist, error)
	DBGetStocksCount(string) (int, error)
	DBGetStocksPaginated(string, int, int) ([]Stock, error)
	GetPriceStatus() (string, error)
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Close() error
}

type FakeRedis struct {
}

type Connectors struct {
	// add mongodb connector here
	session Session
	http    *http.Client
	redis   FakeRedis
	name    string
}

// fake redis Get
func (r *Connectors) Get(key string) (string, error) {
	return store[key], nil
}

// fake redis Set
func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	store[key] = value
	return string(expr), nil
}

// fake redis Close
func (r *Connectors) Close() error {
	store = nil
	return nil
}

// fake mongo section

// DataLayer is an interface to access to the database struct.
type DataLayer interface {
	C(name string) Collection
}

// Collection is an interface to access to the collection struct.
type Collection interface {
	Find(query interface{}) Query
	FindId(query interface{}) Query
	Count() (n int, err error)
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	Update(selector interface{}, update interface{}) error
	// GetMyDocuments() ([]interface{}, error)
	EnsureIndex(index interface{}) error
}

type Query interface {
	All(result interface{}) error
	One(result interface{}) error
	Sort(field interface{}) Query
	Count() (int, error)
	Limit(val int) Query
	Skip(val int) Query
	Iter() FakeIter
	//Iter() FakeIter
}

type Iter interface {
	Next(data interface{}) bool
	Err() error
	Close()
}

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	Close()
	Clone() Session
}

// FakeSession satisfies Session and act as a mock of *mgo.session.
type FakeSession struct{}

// NewFakeSession mock NewSession.
func NewFakeSession() Session {
	return FakeSession{}
}

func (fs FakeSession) Clone() Session {
	return FakeSession{}
}

// Close fakes mgo.Session.Close().
func (fs FakeSession) Close() {}

// DB fakes mgo.Session.DB().
func (fs FakeSession) DB(name string) DataLayer {
	fakeDatabase := FakeDatabase{}
	return fakeDatabase
}

// FakeDatabase satisfies DataLayer and act as a mock.
type FakeDatabase struct{}

// C fakes mgo.Database(name).Collection(name).
func (db FakeDatabase) C(name string) Collection {
	return FakeCollection{Name: name}
}

// FakeCollection satisfies Collection and act as a mock.
type FakeCollection struct {
	Name string
}

// Find fake.
func (fc FakeCollection) Find(query interface{}) Query {
	fq := FakeQuery{Name: fc.Name}
	return fq
}

// Find fake.
func (fc FakeCollection) FindId(query interface{}) Query {
	fq := FakeQuery{Name: fc.Name}
	return fq
}

// Count fake.
func (fc FakeCollection) Count() (n int, err error) {
	return 10, nil
}

// Insert fake.
func (fc FakeCollection) Insert(docs ...interface{}) error {
	return nil
}

// Remove fake.
func (fc FakeCollection) Remove(selector interface{}) error {
	return nil
}

// Update fake.
func (fc FakeCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

// EnsurIndex fake.
func (fc FakeCollection) EnsureIndex(index interface{}) error {
	return nil
}

// GetMyDocuments fake.
func (fc FakeCollection) GetMyDocuments(file string) ([]interface{}, error) {
	var documents []interface{}
	content, _ := ioutil.ReadFile(file)
	json.Unmarshal(content, &documents)

	return documents, nil
}

// FakeQuery satisfies Query and act as a mock.
type FakeQuery struct {
	Name string
}

// All fake.
func (fq FakeQuery) All(result interface{}) error {
	return nil
}

// One fake.
func (fq FakeQuery) One(result interface{}) error {
	if reflect.TypeOf(result).String() == "*main.Stock" {
		*result.(*Stock) = Stock{UID: bson.ObjectIdHex("5cc042307ccc69ada893144c"), PublicationId: 123, AffiliateId: "SBR-01", RefId: 1, Symbol: "TST", Name: "TestSymbol", Buy: 2.0, Stop: 1.0, Last: 3.0, Change: 23.0, Recommendation: "Sell", Status: 1}
	}
	if reflect.TypeOf(result).String() == "*main.Watchlist" {
		*result.(*Watchlist) = Watchlist{UID: bson.ObjectIdHex("5cc042307ccc69ada893144c"), CustomerId: "123", Stocks: []string{}}
	}
	return nil
}

// Distinct fake.
func (fq FakeQuery) Distinct(field string, result interface{}) error {
	return nil
}

// Distinct fake.
func (fq FakeQuery) Count() (int, error) {
	return 10, nil
}

type FakeIter struct {
}

func (f FakeQuery) Iter() FakeIter {
	fio := FakeIter{}
	return fio
}

// Sort fake.
func (fq FakeQuery) Sort(field interface{}) Query {
	return fq
}

// Limit fake.
func (fq FakeQuery) Limit(val int) Query {
	return fq
}

func (fq FakeQuery) Skip(val int) Query {
	return fq
}

func (fi FakeIter) Next(data interface{}) bool {
	counter++
	if reflect.TypeOf(data).String() == "*main.Publication" {
		*data.(*Publication) = Publication{Id: 1, Name: "Test1", AffiliateId: "SBR-01"}
		if counter > 2 {
			counter = 0
			return false
		} else {
			return true
		}
	} else if reflect.TypeOf(data).String() == "*main.Stock" {
		*data.(*Stock) = Stock{UID: bson.ObjectIdHex("5cc042307ccc69ada893144c"), PublicationId: 123, AffiliateId: "SBR-01", RefId: 1, Symbol: "TST", Name: "TestSymbol", Buy: 2.0, Stop: 1.0, Last: 3.0, Change: 23.0, Recommendation: "Sell", Status: 1}
		if counter > 2 {
			counter = 0
			return false
		} else {
			return true
		}
	}
	return false
}

func (fi FakeIter) Err() error {
	return nil
}

func (fi FakeIter) Close() {
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestClients(data string, code int) Clients {

	// read the config
	//config, _ = Init("config.json")
	logger.Level = "info"

	// initialise our store (cache)
	store = make(map[string]string)
	// in initialise the store
	store[DBUPDATESTOCKCURRENTPRICE] = "test"

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	file, _ := ioutil.ReadFile(data)
	logger.Debug(fmt.Sprintf("File %s with data %s", data, string(file)))
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(file))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	redisclient := FakeRedis{}
	mongoclient := NewFakeSession()
	conns := &Connectors{redis: redisclient, session: mongoclient, http: httpclient, name: "test"}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestAll(t *testing.T) {

	// create anonymous struct
	tests := []struct {
		Name     string
		Payload  string
		Handler  string
		FileName string
		want     bool
		errorMsg string
	}{
		{
			"DBSetup should pass",
			"[{\"id\": \"SBR-01\", \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}]",
			"DBSetup",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBSetup should fail",
			`[{"test":"]`,
			"DBSetup",
			"tests/tss.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBIndex should pass",
			"[{\"id\": \"SBR-01\", \"name\":\"Test\",\"token\": \"sdasdsafsfdgdfgf\"}]",
			"DBIndex",
			"tests/tss.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBMigrate should pass",
			"{\"id\": \"SBR-01\", \"affiliate\":\"Test\"}",
			"DBMigrate",
			"tests/publication.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBMigrate should fail",
			`{"test":"`,
			"DBMigrate",
			"tests/publication.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateAffiliateSpecific should pass",
			"{\"id\": \"SBR-01\", \"affiliate\":\"Test\"}",
			"DBUpdateAffiliateSpecific",
			"tests/tss.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateAffiliateSpecific should fail",
			`{"test":"`,
			"DBUpdateAffiliateSpecific",
			"tests/alphavantage.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateStockCurrentPrice should pass",
			"{\"id\": \"SBR-01\", \"affiliate\":\"Test\"}",
			"DBUpdateStockCurrentPrice",
			"tests/alphavantage.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateStock should pass",
			"{\"_id\": \"5cc042307ccc69ada893144c\"}",
			"DBUpdateStock",
			"tests/tss.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateStock should fail",
			`{"test":"`,
			"DBUpdateStock",
			"tests/tss.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetAffiliates should pass",
			"{\"id\": \"SBR-01\", \"affiliate\":\"Test\"}",
			"DBGetAffiliates",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetPublications should pass",
			"{\"id\": \"SBR-01\", \"affiliate\":\"Test\"}",
			"DBGetPublications",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetStocks should pass",
			"5cc042307ccc69ada893144c",
			"DBGetStocks",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetStocks should pass",
			"0",
			"DBGetStocks",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBUpdateWatchlist should pass",
			"{\"_id\": \"5cc042307ccc69ada893144c\", \"affiliate\":\"Test\"}",
			"DBUpdateWatchlist",
			"tests/publication.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetWatchlist should pass",
			"1",
			"DBGetWatchlist",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"GetPriceStatus should pass",
			"",
			"GetPriceStatus",
			"tests/publication.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetStocksCount should pass",
			"1",
			"DBGetStocksCount",
			"tests/publication.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"DBGetStocksPaginated should pass",
			"1",
			"DBGetStocksPaginated",
			"tests/publication.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
	}

	var err error
	for _, tt := range tests {
		logger.Info(fmt.Sprintf("Executing test : %s \n", tt.Name))
		switch tt.Handler {
		case "DBSetup":
			logger.Debug(fp("Payload", tt.Payload))
			connectors = NewTestClients(tt.FileName, 200)
			err = connectors.DBSetup([]byte(tt.Payload))
		case "DBIndex":
			connectors = NewTestClients(tt.FileName, 200)
			err = connectors.DBIndex()
		case "DBMigrate":
			connectors = NewTestClients(tt.FileName, 200)
			err = connectors.DBMigrate([]byte(tt.Payload))
		case "DBUpdateAffiliateSpecific":
			connectors = NewTestClients(tt.FileName, 200)
			err = connectors.DBUpdateAffiliateSpecific([]byte(tt.Payload))
		case "DBUpdateStockCurrentPrice":
			connectors = NewTestClients(tt.FileName, 200)
			err = connectors.DBUpdateStockCurrentPrice()
		case "DBUpdateStock":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBUpdateStock([]byte(tt.Payload))
		case "DBGetAffiliates":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetAffiliates()
		case "DBGetPublications":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetPublications(tt.Payload)
		case "DBGetStocks":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetStocks(tt.Payload, true)
		case "DBUpdateWatchlist":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBUpdateWatchlist([]byte(tt.Payload))
		case "DBGetWatchlist":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetWatchlist(tt.Payload)
		case "GetPriceStatus":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.GetPriceStatus()
		case "DBGetStocksCount":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetStocksCount(tt.Payload)
		case "DBGetStocksPaginated":
			connectors = NewTestClients(tt.FileName, 200)
			_, err = connectors.DBGetStocksPaginated(tt.Payload, 0, 10)
		}

		if !tt.want {
			if err != nil {
				t.Errorf(fmt.Sprintf(tt.errorMsg, tt.Handler, err, nil))
			}
		} else {
			if err == nil {
				t.Errorf(fmt.Sprintf(tt.errorMsg, tt.Handler, "nil", "error"))
			}
		}
	}
}
