package im

import (
	"errors"
	"goim/helpers"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

//add a transfer wrapper
func KeyToUserID(key string) string {
	tmp := ""
	md5 := ""
	strs := strings.Split(key, "\n")
	for _, str := range strs {
		if strings.Contains(str, "--") {
			continue
		}
		tmp += str
		if len(str) > 0 {
			md5 += string(str[0])
		}
	}
	md5 += helpers.MD5([]byte(tmp))
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

//key,value key also be id
type Contact struct {
	//public key
	ID       string
	Key      string
	Nickname string
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

func (p *UserInfo) GetAllContacts() ([]*Contact, error) {
	contacts := []*Contact{}
	err := OpenDB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(p.Key))

		return b.ForEach(func(_, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			var c Contact
			ex := helpers.UnMarshal(v, &c)
			if ex == nil {
				contacts = append(contacts, &c)
			}
			return ex
		})
	})
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

//parameter is public key
func (p *UserInfo) GetContactByID(id string) (*Contact, error) {
	var content []byte
	OpenDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(p.Key))
		if b != nil {
			content = b.Get([]byte(id))
		}
		return nil
	})
	if content != nil {
		var c Contact
		e := helpers.UnMarshal(content, &c)
		if e != nil {
			return nil, e
		}
		return &c, nil
	}
	return nil, errors.New("not exists")
}

//parameter is public key
func (p *UserInfo) SetContact(key string) error {
	id := KeyToUserID(key)
	u := Contact{
		ID:  id,
		Key: key,
	}

	bs, err := helpers.Marshal(&u)
	if err != nil {
		return err
	}
	return OpenDB().Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(p.Key))
		if e != nil {
			return e
		}
		return b.Put([]byte(id), bs)
	})
}

//remove contact.
func (p *UserInfo) RemoveContact(key string) error {
	id := KeyToUserID(key)
	return OpenDB().Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(p.Key))
		if e != nil {
			return e
		}
		return b.Delete([]byte(id))
		//return b.Put([]byte(id), bs)
	})
}

type FriendRequest struct {
	FromUserID   string
	HelloMsg     string
	ModifiedTime time.Time
	IsApproved   bool
}

func getFriendRequestBucketName(prifix string) []byte {
	name := prifix + ".requests"
	return []byte(name)
}

//friend requests
//buckets should be key+"requests"
func (p *UserInfo) GetAllRequests() ([]*FriendRequest, error) {
	reqs := []*FriendRequest{}
	err := OpenDB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(getFriendRequestBucketName(p.Key))

		return b.ForEach(func(_, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			var c FriendRequest
			ex := helpers.UnMarshal(v, &c)
			if ex == nil {
				reqs = append(reqs, &c)
			}
			return ex
		})
	})
	if err != nil {
		return nil, err
	}
	return reqs, nil
}

func (p *UserInfo) AddNewRequest(key string, helloMsg string) error {
	return p.SetRequest(key, helloMsg, false)
}

func (p *UserInfo) AddNewRequestByID(id string, helloMsg string) error {
	return p.SetRequestByID(id, helloMsg, false)
}

func (p *UserInfo) SetRequest(key string, helloMsg string, approve bool) error {
	id := KeyToUserID(key)
	return p.SetRequestByID(id, helloMsg, approve)
}

func (p *UserInfo) SetRequestByID(id string, helloMsg string, approve bool) error {
	u := FriendRequest{
		FromUserID:   id,
		HelloMsg:     helloMsg,
		ModifiedTime: time.Now(),
		IsApproved:   approve,
	}

	bs, err := helpers.Marshal(&u)
	if err != nil {
		return err
	}
	return OpenDB().Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists(getFriendRequestBucketName(p.Key))
		if e != nil {
			return e
		}
		return b.Put([]byte(id), bs)
	})
}
func (p *UserInfo) GetRequest(key string) (*FriendRequest, error) {
	id := KeyToUserID(key)
	return p.GetRequestByID(id)
}

func (p *UserInfo) GetRequestByID(id string) (*FriendRequest, error) {
	var req FriendRequest
	exists := true
	err := OpenDB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(getFriendRequestBucketName(p.Key))

		v := b.Get([]byte(id))
		if v == nil || len(v) == 0 {
			exists = false
			return nil
		}
		return helpers.UnMarshal(v, &req)
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return &req, nil
	}
	return nil, nil
}

func (p *UserInfo) IsRequestExists(key string) (bool, error) {
	req, e := p.GetRequest(key)
	if e != nil {
		return false, e
	}
	if req == nil {
		return false, nil
	} else {
		return true, nil
	}
}

//remove contact.
func (p *UserInfo) ApproveRequestByID(id string) error {
	req, e := p.GetRequestByID(id)
	if e != nil {
		return e
	}
	if req == nil {
		return errors.New("request not exists")
	}
	if req.IsApproved {
		return nil
	}
	return p.SetRequestByID(id, req.HelloMsg, true)
}
