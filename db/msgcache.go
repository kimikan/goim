package db

import (
	"errors"
	"goim/helpers"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	CacheFile = "msgs.db"
)

var dbContext *bolt.DB

func openMsgsDB() *bolt.DB {
	if dbContext == nil {
		db, e := bolt.Open(CacheFile, 0600, nil)
		if e != nil {
			log.Fatal(e)
		}
		dbContext = db
	}
	return dbContext
}

func closeDB() {
	if dbContext != nil {
		dbContext.Close()
		dbContext = nil
	}
}

type TextItem struct {
	Content      string
	FromUserID   string
	ModifiedTime time.Time
}

func GetAllMsgs(id string) ([]*TextItem, error) {
	msgs := []*TextItem{}
	err := openMsgsDB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b, e := tx.CreateBucketIfNotExists([]byte(id))
		if e != nil {
			return e
		}

		return b.ForEach(func(_, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			var c TextItem
			ex := helpers.UnMarshal(v, &c)
			if ex == nil {
				msgs = append(msgs, &c)
			}
			return ex
		})
	})
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

//parameter is public key
func AddMessage(id string, msg *TextItem) error {
	if msg == nil {
		return errors.New("invalid parameter")
	}

	bs, err := helpers.Marshal(msg)
	if err != nil {
		return err
	}
	return openMsgsDB().Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(id))
		if e != nil {
			return e
		}
		k, e2 := msg.ModifiedTime.MarshalBinary()
		if e2 != nil {
			return e2
		}
		return b.Put(k, bs)
	})
}

//remove contact.
func RemoveAllMsgs(id string) error {
	return openMsgsDB().Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(id))
		//return b.Put([]byte(id), bs)
	})
}
