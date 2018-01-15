package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

var (
	mainDB *mgo.Database
	urlDB  = "mongody"
	defDB  = "vault"
)

// InitMainDb : Initialize db connection, should only be called once
func InitMainDb() {
	sess, err := mgo.Dial(urlDB)
	if err != nil {
		log.Fatal(err)
	}
	mainDB = sess.DB(defDB)
}

// CloseMainDB : Close main db session
func CloseMainDB() {
	if mainDB != nil {
		mainDB.Session.Close()
	}
}

// GetNewDB : Dials a new connection, Used to run test
func GetNewDB(dbURL string, dbName string) *mgo.Database {
	sess, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	mDB := sess.DB(dbName)
	return mDB
}

// CopyMainDB : Copy session
func CopyMainDB() *mgo.Database {
	return mainDB.Session.Copy().DB(defDB)
}

// CloseDbSession : Close session
func CloseDbSession(s *mgo.Session) {
	s.Close()
}
