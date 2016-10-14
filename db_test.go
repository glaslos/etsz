package etsz

import (
	"testing"
	"time"

	"github.com/dgryski/go-tsz"
	"github.com/stretchr/testify/assert"
)

func TestGetDBInstance(t *testing.T) {
	edb := New()
	db1 := edb.getDB("test")
	db2 := edb.getDB("test")
	assert.Equal(t, db1, db2, "second call to getDB() should return db1")
}

func TestInsert(t *testing.T) {
	edb := New()
	edb.Insert(1.0, "")
	edb.Insert(1.0, "test")
}

func TestRead(t *testing.T) {
	edb := New()
	edb.Read("test")
	edb.Read("")
}

func TestInsertRead(t *testing.T) {
	edb := New()
	edb.Insert(1.0, "test")
	_ = edb.Read("test")
}

func TestDeleteShard(t *testing.T) {
	edb := New()
	_ = edb.getDB("test")
	// Adding a "random" shard before now
	edb.DBList["test"]["20060102T15"] = tsz.New(uint32(time.Now().Unix()))
	// Insert deletes old shards
	edb.Insert(1.0, "test")
}
