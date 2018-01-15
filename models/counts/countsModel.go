package counts

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Count :
type Count struct {
	Name     string `json:"name" bson:"name"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

// Counts collection
const countsCL = "counts"

// CountConn : Count database connection
type CountConn struct {
	db *mgo.Database
	cl *mgo.Collection
}

// NewCountConn : Get new connection for "counts"
func NewCountConn(idb *mgo.Database) *CountConn {
	return &CountConn{
		db: idb,
		cl: idb.C(countsCL),
	}
}

// AddCount : Increment Count
func (cConn *CountConn) AddCount(cname string) error {
	// Ensure unique Name
	err := cConn.cl.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	if err != nil {
		return fmt.Errorf("AddCount: couldn't add name index: %v", err)
	}

	// Update or Insert value
	_, err = cConn.cl.Upsert(bson.M{"name": cname},
		bson.M{
			"$inc": bson.M{"quantity": 1},
			"$set": bson.M{"name": cname},
		})
	if err != nil {
		return fmt.Errorf("AddCount: upsert failed: %v", err)
	}

	return nil
}

// GetCount : Returns count object given name
func (cConn *CountConn) GetCount(cname string) (*Count, error) {
	c := &Count{}
	if err := cConn.cl.Find(bson.M{"name": cname}).One(c); err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteCount : Delete counter given name
func (cConn *CountConn) DeleteCount(cname string) error {
	return cConn.cl.Remove(bson.M{"name": cname})
}
