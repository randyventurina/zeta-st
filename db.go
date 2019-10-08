package main

import (
	"github.com/syndtr/goleveldb/leveldb"
)


//Db wraps leveldb related operations
type Db struct {
}

//NewDb connects to leveldb specified in application configuration
func NewDb() (*leveldb.DB, error) {
	return leveldb.OpenFile("db", nil)
}
