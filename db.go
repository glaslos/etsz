package etsz

import (
	"sort"
	"time"

	"github.com/dgryski/go-tsz"
)

type entry struct {
	Date  int32   `json:"date"`
	Value float64 `json:"value"`
}

type data struct {
	Date  int32   `json:"date"`
	Value float64 `json:"value"`
}

type shard struct {
	Range string `json:"range"`
	Data  []data `json:"data"`
}

// Result from reading the time series
type Result struct {
	Name   string  `json:"name"`
	Shards []shard `json:"shards"`
}

// EDB is an embedded time series database
type EDB struct {
	// DBList is the list of databases instanciated
	DBList map[string]map[string]*tsz.Series
}

// New creates a new EDB instance
func New() EDB {
	return EDB{DBList: make(map[string]map[string]*tsz.Series)}
}

func (edb *EDB) getDB(name string) *tsz.Series {
	now := time.Now()
	// Shard by minute
	nowString := now.Format("20060102T15")
	if db, ok := edb.DBList[name][nowString]; ok {
		return db
	}
	if len(edb.DBList[name]) == 0 {
		edb.DBList[name] = make(map[string]*tsz.Series)
	}
	edb.DBList[name][nowString] = tsz.New(uint32(now.Unix()))
	return edb.DBList[name][nowString]
}

// Insert data into time series
func (edb *EDB) Insert(v float64, databaseName string) bool {
	if databaseName == "" {
		databaseName = "default"
	}
	now := time.Now()
	edb.getDB(databaseName).Push(uint32(now.Unix()), v)
	currentDayString := now.Format("20060102")
	for shardName := range edb.DBList[databaseName] {
		if shardName[:8] != currentDayString {
			delete(edb.DBList[databaseName], shardName)
		}
	}
	return true
}

func aggr(shardDB *tsz.Series) map[time.Time]float64 {
	aggregated := map[time.Time]float64{}
	it := shardDB.Iter()
	for it.Next() {
		tt, vv := it.Values()
		aggregated[time.Unix(int64(tt), 0)] += vv
	}
	return aggregated
}

// Read data from time series
func (edb *EDB) Read(databaseName string) Result {
	rr := Result{Name: databaseName, Shards: []shard{}}
	db := edb.DBList[databaseName]
	for shardName, shardDB := range db {
		dataShard := shard{Range: shardName, Data: []data{}}
		// Aggregate the data by
		aggregated := aggr(shardDB)
		it := shardDB.Iter()
		for it.Next() {
			tt, vv := it.Values()
			aggregated[time.Unix(int64(tt), 0)] += vv
		}
		for tt, vv := range aggregated {
			if tt.After(time.Now().Add(-60 * time.Minute)) {
				dataShard.Data = append(dataShard.Data, data{int32(tt.Unix()), vv})
			}
		}
		rr.Shards = append(rr.Shards, dataShard)
	}
	return rr
}

// ReadAll data from all time series
func (edb *EDB) ReadAll() []Result {
	r := []Result{}
	for dbName := range edb.DBList {
		rr := edb.Read(dbName)
		r = append(r, rr)
	}
	return r
}

// ReadInt reads data from time series and returns an integer array
// TODO (glaslos): This is some ugly shit...
func (edb *EDB) ReadInt(databaseName string) []int {
	db := edb.getDB(databaseName)
	aggregated := aggr(db)
	// To store the keys in slice in sorted order
	var keys []int
	for k := range aggregated {
		keys = append(keys, int(k.Unix()))
	}
	sort.Ints(keys)

	vals := []int{}
	for _, k := range keys {
		vals = append(vals, int(aggregated[time.Unix(int64(k), 0)]))
	}
	return vals
}
