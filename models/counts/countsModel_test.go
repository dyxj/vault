package counts_test

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
	"vault/db"
	"vault/models/counts"
)

const urlDB = "mongodb://localhost:27017/vault"
const defDB = "vault"

func TestCount(t *testing.T) {
	mdb := db.GetNewDB(urlDB, defDB)
	defer mdb.Session.Close()
	cConn := counts.NewCountConn(mdb)
	cname := "randomtest"

	// Remove existing records
	err := cConn.DeleteCount(cname)
	if err != nil && err != mgo.ErrNotFound {
		t.Fatalf("Failed at Delete: %v", err)
	}

	// Check Get there are no records
	_, err = cConn.GetCount(cname)
	if err != mgo.ErrNotFound {
		t.Fatalf("Failed at Get1, should be error not found: %v", err)
	}

	// Add count
	err = cConn.AddCount(cname)
	if err != nil {
		t.Fatalf("Failed at Add: %v", err)
	}

	// Check Quantity should be 1
	c, err := cConn.GetCount(cname)
	if err != nil {
		t.Fatalf("Failed at Get2: %v", err)
	}
	if c.Quantity != 1 {
		t.Fatal("Quantity should be 1")
	}

	// Add count
	err = cConn.AddCount(cname)
	if err != nil {
		t.Fatalf("Failed at Add: %v", err)
	}

	// Check Quantity should be 2
	c, err = cConn.GetCount(cname)
	if err != nil {
		t.Fatalf("Failed at Get3: %v", err)
	}
	if c.Quantity != 2 {
		t.Fatal("Quantity should be 2")
	}

	// Remove existing records
	err = cConn.DeleteCount(cname)
	if err != nil && err != mgo.ErrNotFound {
		t.Fatalf("Failed at Delete: %v", err)
	}

	// Check Get there are no records
	_, err = cConn.GetCount(cname)
	if err != mgo.ErrNotFound {
		t.Fatalf("Failed at Get4, should be error not found: %v", err)
	}
}
