package im

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

//a concreate storage implementation using leveldb
//the definition of the storemanager
type StoreManager struct {
	sync.RWMutex
	UserDb    *leveldb.DB
	UserIter  iterator.Iterator
	GroupDb   *leveldb.DB
	GroupIter iterator.Iterator
	Users  map[[16]byte]*User
	Groups map[[16]byte]*Group
}

//the creator ops
func NewStore() *StoreManager {
	p := &StoreManager{}

	db, err := leveldb.OpenFile("userdb", nil)
	if err != nil {
		p.UserDb = db
	}

	db, err = leveldb.OpenFile("groupdb", nil)
	if err != nil {
		p.GroupDb = db
	}
	p.Users = make(map[[16]byte]*User)
	p.Groups = make(map[[16]byte]*Group)
	return p
}

func (p *StoreManager) Close() {
	if p.UserDb != nil {
		p.UserDb.Close()
	}
	if p.GroupDb != nil {
		p.GroupDb.Close()
	}
}

//insert a new user
func (p *StoreManager) AddUser(u *User) error {
	bytes, err := u.ToBytes()
	if err != nil {
		return nil
	}

	err = p.UserDb.Put(u.Name[:], bytes, nil)
	if err != nil {
		return err
	}

	return nil
}

//not implemented
func (p *StoreManager) FirstUser() *User {
	p.UserIter = p.UserDb.NewIterator(nil, nil)

	return p.NextUser()
}

func (p *StoreManager) NextUser() *User {
	if p.UserIter != nil {
		if p.UserIter.Next() {
			u := p.UserIter.Value()
			user, err := ParseUser(u)
			if err != nil {
				return user
			}
		}
	}
	return nil
}

func (p *StoreManager) GetUser(id []byte) *User {
	var id2 [16]byte
	copy(id2[:], id)
	if v, ok := p.Users[id2]; ok {
		return v
	}
	bytes, err := p.UserDb.Get(id, nil)
	if err != nil {
	return nil
}

	user, err := ParseUser(bytes)
	if err != nil {
	return nil
	}
	p.Users[id2] = user
	return user
}

func (p *StoreManager) DeleteUser(id []byte) error {
	var id2 [16]byte
	copy(id2[:], id)
	delete(p.Users, id2)
	return p.UserDb.Delete(id, nil)
}
//insert a new group
func (p *StoreManager) AddGroup(g *Group) error {
	bytes, err := g.ToBytes()
	if err != nil {
		return nil
	}
	err = p.GroupDb.Put(g.Name[:], bytes, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *StoreManager) FirstGroup() *Group {
	p.GroupIter = p.GroupDb.NewIterator(nil, nil)
	return p.NextGroup()
}

func (p *StoreManager) NextGroup() *Group {
	if p.GroupIter != nil {
		if p.GroupIter.Next() {
			u := p.GroupIter.Value()
			group, err := ParseGroup(p, u)
			if err != nil {
				return group
			}
		}
	}
	return nil
}

func (p *StoreManager) GetGroup(id []byte) *Group {
	var id2 [16]byte
	copy(id2[:], id)
	if v, ok := p.Groups[id2]; ok {
		return v
	}
	bytes, err := p.GroupDb.Get(id, nil)
	if err != nil {
	return nil
}

	group, err := ParseGroup(p, bytes)
	if err != nil {
	return nil
	}
	p.Groups[id2] = group
	return group
}
func (p *StoreManager) DeleteGroup(id []byte) error {
	var id2 [16]byte
	copy(id2[:], id)
	delete(p.Groups, id2)
	return p.GroupDb.Delete(id, nil)
}
