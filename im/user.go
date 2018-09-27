package im

import (
	"errors"
	"goim/helpers"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

//add a transfer wrapper
func KeyToUserID(key string) string {
	md5 := helpers.MD5([]byte(key))
	strs := strings.Split(key, "\n")
	for _, str := range strs {
		if strings.Contains(str, "--") {
			continue
		}
		if len(str) > 0 {
			md5 += string(str[0])
		}
	}
	return md5
}

const (
	DbFile            = "im.db"
	bucketUserProfile = "user_profile_table"
)

var dbContext *bolt.DB

func OpenDB() *bolt.DB {
	if dbContext == nil {
		db, e := bolt.Open(DbFile, 0600, nil)
		if e != nil {
			log.Fatal(e)
		}
		dbContext = db
	}
	return dbContext
}

func CloseDB() {
	if dbContext != nil {
		dbContext.Close()
		dbContext = nil
	}
}

//UserID=>UserInfo
//UserKey=>UserProfile
type UserInfo struct {
	//public key
	Key         string
	DisplayName string
	Description string
	Avatar      []byte
}

func GetUserInfoByPubKey(key string) (*UserInfo, error) {
	id := KeyToUserID(key)
	return GetUserInfoByID(id)
}

func GetUserInfoByID(id string) (*UserInfo, error) {
	var content []byte
	OpenDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUserProfile))
		if b != nil {
			content = b.Get([]byte(id))
		}
		return nil
	})
	if content != nil {
		var user UserInfo
		e := helpers.UnMarshal(content, &user)
		if e != nil {
			return nil, e
		}
		return &user, nil
	}
	return nil, errors.New("not exists")
}

func SetUserInfo(u *UserInfo) error {
	if u == nil {
		return errors.New("SetUserInfo: invalid parameter")
	}
	id := KeyToUserID(u.Key)
	return SetUserInfoByID(id, u)
}

//set userinfo by id, should.
func SetUserInfoByID(id string, u *UserInfo) error {
	if u == nil {
		return errors.New("invalid user")
	}
	bs, err := helpers.Marshal(u)
	if err != nil {
		return err
	}
	return OpenDB().Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(bucketUserProfile))
		if e != nil {
			return e
		}
		return b.Put([]byte(id), bs)
	})
}
