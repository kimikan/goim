package im

import (
	"github.com/google/uuid"
)

type Hashable interface {
	HashKey() interface{}
}

type Store interface {
	//insert a new user
	AddUser(u *User) error
	FirstUser() *User
	NextUser() *User
	GetUser(id *uuid.UUID) *User
	DeleteUser(id *uuid.UUID) error

	//insert a new group
	AddGroup(g *Group) error
	FirstGroup() *Group
	NextGroup() *Group
	GetGroup(id *uuid.UUID) *Group
	DeleteGroup(id *uuid.UUID) error
}
