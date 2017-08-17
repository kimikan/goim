package im

/* hashkey */

type Hashable interface {
	HashKey() interface{}
}

//the storage provider
//public a lot of virtual operations
//to do something
type Store interface {
	//insert a new user
	AddUser(u *User) error
	FirstUser() *User
	NextUser() *User
	GetUser(id []byte) *User
	DeleteUser(id []byte) error

	//insert a new group
	AddGroup(g *Group) error
	FirstGroup() *Group
	NextGroup() *Group
	GetGroup(id []byte) *Group
	DeleteGroup(id []byte) error
}
