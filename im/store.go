package im

import (
	"sync"

	"github.com/golang/leveldb"
	"github.com/golang/leveldb/db"
	"github.com/google/uuid"
)

//the definition of the storemanager
type StoreManager struct {
	sync.RWMutex
	Database *leveldb.DB
}

//the creator
func NewStore() *StoreManager {
	p := &StoreManager{}

	db, err := leveldb.Open("db", &db.Options{})
	if err != nil {
		p.Database = db
	}

	return p
}

func (p *StoreManager) Close() {
	if p.Database != nil {
		p.Database.Close()
	}
}

//insert a new user
func (p *StoreManager) AddUser(u *User) error {
	buf, err := u.ID.MarshalBinary()
	if err != nil {

	}
	//p.Database.Set(u.ID, u.
	return nil
}

func (p *StoreManager) FirstUser() *User {
	return nil
}

func (p *StoreManager) NextUser() *User {
	return nil
}

func (p *StoreManager) GetUser(id *uuid.UUID) *User {
	return nil
}

func (p *StoreManager) DeleteUser(id *uuid.UUID) error {
	return nil
}

//insert a new group
func (p *StoreManager) AddGroup(g *Group) error {
	return nil
}

func (p *StoreManager) FirstGroup() *Group {
	return nil
}

func (p *StoreManager) NextGroup() *Group {
	return nil
}

func (p *StoreManager) GetGroup(id *uuid.UUID) *Group {
	return nil
}

func (p *StoreManager) DeleteGroup(id *uuid.UUID) error {
	return nil
}
