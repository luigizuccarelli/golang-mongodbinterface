package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MSGFORMAT                 string = "function %s : %v\n"
	DBSETUP                   string = "DBSetup"
	DBMIGRATE                 string = "DBMigrate"
	DBUPDATEAFFILIATESPECIFIC string = "DBUpdateAffiliateSpecific"
	DBUPDATESTOCKCURRENTPRICE string = "DBUpdateStockCurrentPrice"
	DBUPDATESTOCK             string = "DBUpdateStock"
	DBGETAFFILIATES           string = "DBGetAffiliates"
	DBGETPUBLICATIONS         string = "DBGetPublications"
	DBGETSTOCKS               string = "DBGetStocks"
	DBGETWATCHLIST            string = "DBGetWatchlist"
	DBWATCHLIST               string = "DBUpdateWatchlist"
	AFFILIATE                 string = "affiliate"
	AFFILIATES                string = "affiliates"
	AFFILIATEID               string = "affiliateid"
	PUBLICATIONS              string = "publications"
	PUBLICATIONID             string = "publicationid"
	WATCHLIST                 string = "watchlist"
	STOCKS                    string = "stocks"
	SYMBOL                    string = "symbol"
	STATUS                    string = "status"
	MERGEDDATA                string = " : merged data"
	DATA                      string = " : data"
	CLONE                     string = "Session clone"
)

var lock = sync.RWMutex{}

func fp(msg string, obj interface{}) string {
	return fmt.Sprintf(MSGFORMAT, msg, obj)
}

// DBSetup this function saves the affiliate data
// It takes in a byte array and returns an error
func (c *Connectors) DBSetup(b []byte) error {
	// This function must be run before DBMigrate
	// initial check TBD
	logger.Trace(DBSETUP)
	var affiliates []Affiliate
	// read the payload in the form of []Affiliate
	// parse input and store to db
	e := json.Unmarshal(b, &affiliates)
	if e != nil {
		logger.Error(fp(DBSETUP, e.Error()))
		return e
	}
	logger.Debug(fp(DBSETUP+" Inserting data", affiliates))
	s := c.session.Clone()
	collection := s.DB(config.MongoDB.DatabaseName).C(AFFILIATES)
	defer s.Close()
	// convert to []interface{} for array insert
	var ui []interface{}
	for _, t := range affiliates {
		ui = append(ui, t)
	}
	e = collection.Insert(ui...)
	if e != nil {
		logger.Error(fp(DBSETUP+" Inserting data", e.Error()))
		return e
	}
	return nil
}

// DBIndex this function indexes the various collections
// It has void paramaters and returns an error
func (c *Connectors) DBIndex() error {

	logger.Trace("DBIndex")
	s := c.session.Clone()
	collection := s.DB(config.MongoDB.DatabaseName).C(AFFILIATES)
	index := mgo.Index{
		Key: []string{"id"},
	}
	defer s.Close()
	err := collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	collection = s.DB(config.MongoDB.DatabaseName).C(PUBLICATIONS)
	index = mgo.Index{
		Key: []string{"id", AFFILIATEID},
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	collection = s.DB(config.MongoDB.DatabaseName).C(STOCKS)
	index = mgo.Index{
		Key: []string{"id", PUBLICATIONID, SYMBOL, AFFILIATEID, STATUS},
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		return err
	}

	collection = s.DB(config.MongoDB.DatabaseName).C(WATCHLIST)
	index = mgo.Index{
		Key: []string{"customerid"},
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		return err
	}

	return nil
}

// DBMigrate this function reads data via the tradesmiths api and stores the structure into mongodb
// It takes in a byte array and returns an error
// The receiver Connectors is used to allow for unit testing
func (c *Connectors) DBMigrate(b []byte) error {

	logger.Trace(DBMIGRATE)

	var affiliate Affiliate
	var publications []Publication
	var stocks []Stock
	var list []Stock
	var keys = make(map[string]bool)
	var j map[string]interface{}

	e := json.Unmarshal(b, &j)
	if e != nil {
		logger.Error(fp(DBMIGRATE, e.Error()))
		return e
	}

	affiliateName := fmt.Sprintf("%s", j[AFFILIATE])
	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	collection := s.DB(config.MongoDB.DatabaseName).C(AFFILIATES)

	// find the affiliate info in DB
	// first find the collection with the given ID
	query := bson.M{"name": affiliateName}
	e = collection.Find(query).One(&affiliate)
	logger.Trace(fp(DBMIGRATE+" : affiliate data", affiliate))
	if e != nil {
		logger.Error(fp(DBMIGRATE+" : finding affiliate", e.Error()))
		return e
	}

	// do the api call to get Publications
	req, err := http.NewRequest("GET", config.Url+"ApiPortfolio/GetAllPortfolios/?ApiKey="+affiliate.Token, nil)
	logger.Info(fp("DBMigrate URL info", config.Url+"ApiPortfolio/GetAllPortfolios/?ApiKey="+affiliate.Token))
	resp, err := c.http.Do(req)
	logger.Info(fmt.Sprintf("Retrieving all publication for affiliate %s %d", affiliate.Name, affiliate.Id))
	if err != nil || resp.StatusCode != 200 {
		logger.Error(fp(DBMIGRATE, err.Error()))
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fp(DBMIGRATE, err.Error()))
		return err
	}

	// convert json to schema
	json.Unmarshal(body, &publications)
	for x, _ := range publications {
		logger.Debug(fmt.Sprintf("Publications info %d", publications[x].Id))
		publications[x].AffiliateId = affiliate.Id
		req, err := http.NewRequest("GET", config.Url+"ApiPosition/GetListPositinsByPortfolioId/?ApiKey="+affiliate.Token+"&portfolioid="+strconv.Itoa(publications[x].Id), nil)
		logger.Debug(fp("DBMigrate URL info", config.Url+"ApiPosition/GetListPositinsByPortfolioId/?ApiKey="+affiliate.Token+"&portfolioid="+strconv.Itoa(publications[x].Id)))
		resp, err := c.http.Do(req)
		logger.Info(fp("DBMigrate retrieving all stocks for publication", publications[x].Name))
		if err != nil || resp.StatusCode != 200 {
			logger.Error(fp("DBMigrate retrieving stock info", err.Error()))
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fp(DBMIGRATE, err.Error()))
			return err
		}
		logger.Debug(fmt.Sprintf("DBMigrate json data from url %s", string(body)))
		json.Unmarshal(body, &stocks)
		for y, _ := range stocks {
			stocks[y].PublicationId = publications[x].Id
			stocks[y].AffiliateId = affiliate.Id
			// check for duplicates , dont add to list if it exists
			if _, value := keys[stocks[y].Symbol]; !value {
				keys[stocks[y].Symbol] = true
				list = append(list, stocks[y])
			}
		}
	}

	logger.Trace(fp("DBMigrate publications info", publications))
	logger.Trace(fp("DBMigrate stocks info", list))

	collection = s.DB(config.MongoDB.DatabaseName).C(PUBLICATIONS)
	// store to DB
	// convert to []interface{} for array insert
	var ui []interface{}
	for _, t := range publications {
		ui = append(ui, t)
	}

	e = collection.Insert(ui...)
	if e != nil {
		logger.Error(fp("DBMigrate inserting publications", e.Error()))
		return e
	}
	collection = s.DB(config.MongoDB.DatabaseName).C(STOCKS)

	defer s.Close()

	var ux []interface{}
	for _, t := range list {
		ux = append(ux, t)
	}

	e = collection.Insert(ux...)
	if e != nil {
		logger.Error(fp("DBMigrate inserting stocks", e.Error()))
		return e
	}

	// all good
	return nil
}

// DBUpdateAffiliateSpecific this function reads data via the tradesmiths api and updates the current mongodb with affiliate specific info
// This is in the form Buy, Stop, Status and Recommendation info
// It takes in a byte array and returns an error
// The receiver Connectors is used to allow for unit testing
func (c *Connectors) DBUpdateAffiliateSpecific(b []byte) error {

	logger.Trace(DBUPDATEAFFILIATESPECIFIC)

	var affiliate Affiliate
	var publications []Publication
	var publication Publication
	var tss []TradeSmithSchema
	var stock Stock
	var keys = make(map[string]bool)
	var j map[string]interface{}

	e := json.Unmarshal(b, &j)
	if e != nil {
		logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, e.Error()))
		return e
	}

	affiliateName := fmt.Sprintf("%s", j[AFFILIATE])
	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	collection := s.DB(config.MongoDB.DatabaseName).C(AFFILIATES)

	// find the affiliate info in DB
	// first find the collection with the given ID
	query := bson.M{"name": affiliateName}
	e = collection.Find(query).One(&affiliate)
	logger.Trace(fp(DBUPDATEAFFILIATESPECIFIC+" : affiliate data", affiliate))
	if e != nil {
		logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, e.Error()))
		return e
	}

	// now get all the Publications
	collection = s.DB(config.MongoDB.DatabaseName).C(PUBLICATIONS)
	// first find the collection with the given ID
	query = bson.M{AFFILIATEID: affiliate.Id}

	// first find the collection with the given ID
	iter := collection.Find(query).Sort("name").Iter()

	for iter.Next(&publication) {
		logger.Trace(fp(DBUPDATEAFFILIATESPECIFIC+" publication data", publication))
		publications = append(publications, publication)
	}
	if iter.Err() != nil {
		logger.Error(fp(DBUPDATEAFFILIATESPECIFIC+" : publication data", iter.Err()))
		iter.Close()
		return iter.Err()
	}
	iter.Close()

	// we iterate through each publication and do a request on the tradesmith api for the publication
	// the json is transformed into a schema and the relevant stock is updated

	// do the api call to get Publications
	for x, _ := range publications {
		req, err := http.NewRequest("GET", config.Url+"ApiPosition/GetAllByPortfolioId/?ApiKey="+affiliate.Token+"&portfolioId="+strconv.Itoa(publications[x].Id), nil)
		logger.Debug(fp("DBUpdateAffiliateSpecific URL info", config.Url+"ApiPosition/GetAllByPortfolioid/?ApiKey="+affiliate.Token+"&portfolioId="+strconv.Itoa(publications[x].Id)))
		resp, err := c.http.Do(req)
		logger.Info(fp("DBUpdateAffiliateSpecific retrieving all positions for publication", publications[x].Id))
		if err != nil || resp.StatusCode != 200 {
			logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, err.Error()))
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, err.Error()))
			return err
		}

		// convert json to schema
		e := json.Unmarshal(body, &tss)
		if e != nil {
			logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, err.Error()))
			return err
		}
		logger.Debug(fp(DBUPDATEAFFILIATESPECIFIC+" stock info from TradeSmiths", tss))

		for y, _ := range tss {
			if keys[tss[y].Symbol] {
				logger.Info(fp(DBUPDATEAFFILIATESPECIFIC+" duplicate stock found no updates will be made", tss[y].Symbol))
			} else {
				if tss[y].Symbol != "" {
					// now to a lookup to the DB for the symbol
					st := s.DB(config.MongoDB.DatabaseName).C(STOCKS)
					logger.Trace(fp("DBUpdateAffiliateSpecific looking up stock", tss[y].Symbol))
					query := bson.M{SYMBOL: tss[y].Symbol}
					// first find the collection with the given ID
					err := st.Find(query).One(&stock)
					if err != nil {
						return err
					}
					logger.Info(fp(DBUPDATEAFFILIATESPECIFIC+" Stocks ", stock))
					// update the fields we are interested in
					stock.Buy = tss[y].Buy
					stock.Stop = tss[y].SubTrades[0].SubstradeSetting.Stop
					// golang does not like % in a string - some cleanup is needed
					stock.Recommendation = strings.Replace(tss[y].Recommendation.Info, "%", " percent", -1)
					stock.Status = tss[y].Status
					stock.CurrencySign = tss[y].CurrencySign

					// update the merged data
					query = bson.M{"_id": bson.ObjectIdHex(stock.UID.Hex())}
					logger.Debug(fp(DBUPDATEAFFILIATESPECIFIC+MERGEDDATA, stock))
					e = st.Update(query, stock)
					if e != nil {
						logger.Error(fp(DBUPDATEAFFILIATESPECIFIC+" : updating", err.Error()))
						return e
					}
					// we keep track odf updated symbols to eliminate duplicates
					keys[stock.Symbol] = true
				} else {
					logger.Error(fp(DBUPDATEAFFILIATESPECIFIC, "Empty stock symbol - please verfiy the tradesmiths api"))
				}
			}
		}
	}
	// all good
	return nil
}

// DBUpdateStockCurrentPrice this function reads stock data from the db and uses the extrenal api to update current stock price
// It takes in a byte array and returns an error
// The receiver Connectors is used to allow for unit testing
func (c *Connectors) DBUpdateStockCurrentPrice() error {

	logger.Trace(DBUPDATESTOCKCURRENTPRICE)

	var stockprice Alphavantage
	var stocks []Stock
	var stock Stock
	var bErr bool = false

	// update redis to indicate stockupdate for prices is pending
	c.Set(DBUPDATESTOCKCURRENTPRICE, "pending", 3600*time.Second)
	lock.Lock()
	defer lock.Unlock()

	go func() {

		// do lookup to get affiliate token on DB
		s := c.session.Clone()
		defer s.Close()
		collection := s.DB(config.MongoDB.DatabaseName).C(STOCKS)

		query := bson.M{"status": 1}
		// find the stocks
		iter := collection.Find(query).Sort(SYMBOL).Iter()

		for iter.Next(&stock) {
			logger.Trace(fp(DBUPDATESTOCKCURRENTPRICE+DATA, stock))
			stocks = append(stocks, stock)
		}
		if iter.Err() != nil {
			logger.Error(fp(DBUPDATESTOCKCURRENTPRICE+DATA, iter.Err()))
			iter.Close()
			return
		}
		iter.Close()

		// iterate through each stock
		for x, _ := range stocks {

			// Get the latest stock data
			req, err := http.NewRequest("GET", config.QuoteUrl+"GLOBAL_QUOTE&symbol="+stocks[x].Symbol+"&apikey="+config.Token, nil)
			resp, err := c.http.Do(req)
			logger.Debug(fp(DBUPDATESTOCKCURRENTPRICE, config.QuoteUrl))
			logger.Info(fp(DBUPDATESTOCKCURRENTPRICE+"Retrieving stock price for ", stocks[x].Symbol))
			if err != nil || resp.StatusCode != 200 {
				// just log the error - this is not a critical error
				logger.Error(fp(DBUPDATESTOCKCURRENTPRICE, err.Error()))
				bErr = true
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(fp(DBUPDATESTOCKCURRENTPRICE, err.Error()))
				bErr = true
			}
			e := json.Unmarshal(body, &stockprice)
			if e != nil {
				logger.Error(fp(DBUPDATESTOCKCURRENTPRICE, e.Error()))
				bErr = true
			}
			logger.Info(fp(DBUPDATESTOCKCURRENTPRICE+"Stock from alphavantage ", stockprice))

			if stockprice == (Alphavantage{}) {
				stocks[x].Status = -1
			} else {
				stocks[x].Last, _ = strconv.ParseFloat(stockprice.GlobalQuote.Price, 64)
				stocks[x].Change, _ = strconv.ParseFloat(stockprice.GlobalQuote.ChangePercent[:len(stockprice.GlobalQuote.ChangePercent)-1], 64)
			}

			query := bson.M{"_id": bson.ObjectIdHex(stocks[x].UID.Hex())}
			logger.Debug(fp(DBUPDATESTOCKCURRENTPRICE+MERGEDDATA, stocks[x]))
			e = collection.Update(query, stocks[x])
			if e != nil {
				logger.Error(fp(DBUPDATESTOCKCURRENTPRICE+MERGEDDATA, e.Error()))
				bErr = true
			}
			stockprice = Alphavantage{}
		}

		//lock.Unlock()
		//defer lock.Unlock()

		if !bErr {
			// update redis to indicate end with success
			c.Set(DBUPDATESTOCKCURRENTPRICE, "OK", 43200*time.Second)
		} else {
			// update redis to indicate end with failure
			c.Set(DBUPDATESTOCKCURRENTPRICE, "KO", 43200*time.Second)
		}
	}()

	return nil
}

// DBUpdateStock - it does what it says :)
// It takes a byte array and returns both the Stock array and error objects
func (c *Connectors) DBUpdateStock(body []byte) ([]Stock, error) {

	var data, existing Stock
	var stocks []Stock

	e := json.Unmarshal(body, &data)
	if e != nil {
		logger.Error(fp(DBUPDATESTOCK+" : reading json", e.Error()))
		return stocks, e
	}

	logger.Debug(fp(DBUPDATESTOCK+DATA, data))

	// session copy
	s := c.session.Clone()
	if s == nil {
		logger.Error(fp(DBUPDATESTOCK, CLONE))
		return stocks, errors.New(CLONE)
	}

	defer s.Close()
	// collection publications
	collection := s.DB(config.MongoDB.DatabaseName).C(STOCKS)

	// check the bson id - the payload must include the id - its not taken from the query string
	f := bson.IsObjectIdHex(data.UID.Hex())
	if !f {
		return stocks, errors.New(DBUPDATESTOCK + " bson ObjectId not valid")
	}

	// first find the collection with the given ID
	err := collection.FindId(data.UID).One(&existing)
	if err != nil {
		return stocks, err
	}
	logger.Debug(fp(DBUPDATESTOCK+" : from database", existing))

	// now merge the 2 structs
	// takes the form (dst,src,mode)
	em := mergo.Merge(&existing, data, mergo.WithOverride)
	if em != nil {
		return stocks, err
	}

	// update the merged structs
	query := bson.M{"_id": bson.ObjectIdHex(data.UID.Hex())}
	logger.Debug(fp(DBUPDATESTOCK+MERGEDDATA, existing))
	e = collection.Update(query, existing)
	if e != nil {
		logger.Error(fp(DBUPDATESTOCK+MERGEDDATA, e.Error()))
		return stocks, e
	}

	stocks = append(stocks, existing)
	// all good
	return stocks, nil
}

// DBGetAffiliates - get a list of all affiliates
// It has void parameters and returns an affiliate schema array
func (c *Connectors) DBGetAffiliates() ([]Affiliate, error) {

	logger.Trace(DBGETAFFILIATES)

	var affiliates []Affiliate
	var data Affiliate

	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(AFFILIATES)
	// first find the collection with the given ID
	iter := collection.Find(nil).Sort("_id").Iter()

	for iter.Next(&data) {
		logger.Trace(fp(DBGETAFFILIATES+DATA, data))
		affiliates = append(affiliates, data)
	}
	if iter.Err() != nil {
		logger.Error(fp(DBGETAFFILIATES+DATA, iter.Err()))
		iter.Close()
		return affiliates, iter.Err()
	}
	iter.Close()

	// all good
	return affiliates, nil
}

// DBGetPublications - get a list of all publications
// It has a string id parameter (publication id) and returns a publication schema array
func (c *Connectors) DBGetPublications(id string) ([]Publication, error) {

	logger.Trace(DBGETPUBLICATIONS)

	var publications []Publication
	var data Publication

	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(PUBLICATIONS)
	// first find the collection with the given ID
	affiliateId, _ := strconv.Atoi(id)
	query := bson.M{AFFILIATEID: affiliateId}

	// first find the collection with the given ID
	iter := collection.Find(query).Sort("name").Iter()

	for iter.Next(&data) {
		logger.Trace(fp(DBGETPUBLICATIONS+DATA, data))
		publications = append(publications, data)
	}
	if iter.Err() != nil {
		logger.Error(fp(DBGETPUBLICATIONS+DATA, iter.Err()))
		iter.Close()
		return publications, iter.Err()
	}
	iter.Close()

	// all good
	return publications, nil
}

// DBGetStocks - get a list of stocks by publication or  affiliate
// It has a string id parameter (publication or affiliate id) , a boolean if set true returns all stocks for an affiliate
func (c *Connectors) DBGetStocks(id string, all bool) ([]Stock, error) {

	logger.Trace(DBGETSTOCKS)

	var stocks []Stock
	var data Stock
	var query bson.M

	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(STOCKS)
	// first find the collection with the given ID
	if !all {
		publicationId, _ := strconv.Atoi(id)
		query = bson.M{PUBLICATIONID: publicationId, "status": 1}
	} else {
		affiliateId, _ := strconv.Atoi(id)
		query = bson.M{AFFILIATEID: affiliateId, "status": 1}
	}

	// first find the collection with the given ID
	iter := collection.Find(query).Sort(SYMBOL).Iter()

	for iter.Next(&data) {
		str := strings.Replace(data.Recommendation, "%", " percent", -1)
		data.Recommendation = str
		logger.Trace(fp(DBGETSTOCKS+DATA, data))
		stocks = append(stocks, data)
	}
	if iter.Err() != nil {
		logger.Error(fp(DBGETSTOCKS+DATA, iter.Err()))
		iter.Close()
		return stocks, iter.Err()
	}
	iter.Close()

	// all good
	return stocks, nil
}

// DBGetStocksCount - get a list of stocks by publication or  affiliate
// It has a string id parameter (publication or affiliate id) , a boolean if set true returns all stocks for an affiliate
func (c *Connectors) DBGetStocksCount(id string) (int, error) {

	logger.Trace(DBGETSTOCKS)

	var query bson.M

	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(STOCKS)
	// first find the collection with the given ID
	affiliateId, _ := strconv.Atoi(id)
	query = bson.M{AFFILIATEID: affiliateId, "status": 1}

	// first find the collection
	result, _ := collection.Find(query).Count()

	// all good
	return result, nil
}

// DBGetStocksPaginated - get a list of stocks by publication or  affiliate
// It has a string id parameter (publication or affiliate id) , a boolean if set true returns all stocks for an affiliate
func (c *Connectors) DBGetStocksPaginated(id string, skip int, limit int) ([]Stock, error) {

	logger.Trace(DBGETSTOCKS + "PAGINATED")

	var stocks []Stock
	var data Stock
	var query bson.M

	// do lookup to get affiliate token on DB
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(STOCKS)
	// first find the collection with the given ID
	affiliateId, _ := strconv.Atoi(id)
	query = bson.M{AFFILIATEID: affiliateId, "status": 1}

	iter := collection.Find(query).Sort(SYMBOL).Skip(skip).Limit(limit).Iter()

	for iter.Next(&data) {
		str := strings.Replace(data.Recommendation, "%", " percent", -1)
		data.Recommendation = str
		logger.Trace(fp(DBGETSTOCKS+DATA, data))
		stocks = append(stocks, data)
	}
	if iter.Err() != nil {
		logger.Error(fp(DBGETSTOCKS+DATA, iter.Err()))
		iter.Close()
		return stocks, iter.Err()
	}
	iter.Close()
	// all good
	return stocks, nil
}

// DBGetWatchList - this function returns a watchlist for a specific client
// It takes a string id parameter and returns a watchlist schema
func (c *Connectors) DBGetWatchlist(id string) (Watchlist, error) {

	logger.Debug(DBGETWATCHLIST + " " + id)

	var data Watchlist
	s := c.session.Clone()
	defer s.Close()
	collection := s.DB(config.MongoDB.DatabaseName).C(WATCHLIST)
	customerId, _ := strconv.Atoi(id)
	query := bson.M{"customerid": customerId}

	// first find the collection with the given ID
	err := collection.Find(query).One(&data)
	if err != nil {
		if strings.Index(err.Error(), "not found") != -1 {
			logger.Warn(fp(DBWATCHLIST+" "+id, err.Error()))
			return data, nil
		} else {
			logger.Error(fp(DBWATCHLIST+" "+id, err.Error()))
			return data, err
		}
	}
	logger.Debug(fp(DBGETWATCHLIST+" : from database", data))
	return data, nil
}

// DBUpdateWatchlist
// It takes a byte array and returns both the ShcemaInterface and error objects
func (c *Connectors) DBUpdateWatchlist(body []byte) (Watchlist, error) {

	var data, existing Watchlist

	e := json.Unmarshal(body, &data)
	if e != nil {
		logger.Error(fp(DBWATCHLIST+" : reading json", e.Error()))
		return data, e
	}

	logger.Debug(fp(DBWATCHLIST+DATA, data))

	// session copy
	s := c.session.Clone()
	if s == nil {
		logger.Error(fp(DBWATCHLIST+" : session", CLONE))
		return data, errors.New(CLONE)
	}

	defer s.Close()
	// collection publications
	collection := s.DB(config.MongoDB.DatabaseName).C("watchlist")

	query := bson.M{"customerid": data.CustomerId}

	// first find the collection with the given ID
	err := collection.Find(query).One(&existing)
	if err != nil {
		// no record found lets insert
		logger.Debug(fp(DBWATCHLIST+" : no record found inserting into database", data))
		// now merge the 2 structs
		// takes the form (dst,src,mode)
		em := mergo.Merge(&existing, data, mergo.WithOverride)
		if em != nil {
			return data, err
		}
		e = collection.Insert(existing)
		if e != nil {
			logger.Error(fp(DBWATCHLIST+" : Inserting watchlist", e.Error()))
			return data, e
		}
	} else {
		// always clear the stocks field
		existing.Stocks = nil
		logger.Debug(fp(DBWATCHLIST+" : record found updating database", existing))

		// now merge the 2 structs
		// takes the form (dst,src,mode)
		em := mergo.Merge(&existing, data, mergo.WithOverride)
		if em != nil {
			return data, err
		}

		// update the merged structs
		// query := bson.M{"_id": bson.ObjectIdHex(data.UID.Hex())}
		logger.Debug(fp(DBWATCHLIST+MERGEDDATA, existing))
		e = collection.Update(query, existing)
		if e != nil {
			logger.Error(fp(DBWATCHLIST+MERGEDDATA, e.Error()))
			return data, e
		}
	}

	// all good
	return existing, nil
}

// GetPriceStatus - get the current price update status update
// It has a void parameter and returns string
func (c *Connectors) GetPriceStatus() (string, error) {
	lock.RLock()
	defer lock.RUnlock()
	val, _ := c.Get(DBUPDATESTOCKCURRENTPRICE)
	return val, nil
}
