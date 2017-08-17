package im

import "bytes"


type Group struct {
	//represent a uuid
	Name [16]byte

	//for privilege approve etc
	Owner *User

	Users []*User
	Messages chan *Message
}

//deserialize the buf to struct
func ParseGroup(store Store, bs []byte) (*Group, error) {
	p := NewGroup("", nil)
	buf := bytes.NewBuffer(bs)

	len, err := buf.Read(p.Name[:])
	if len != 16 || err != nil {
		return nil, err
	}

	tmp := [16]byte{}
	len, err = buf.Read(tmp[:])
	if len != 16 || err != nil {
		return nil, err
	}

	p.Owner = store.GetUser(tmp[:])

	for {
		len, err = buf.Read(tmp[:])
		if len != 16 || err != nil {
			break
		}
		//Find specific user

		user2 := store.GetUser(tmp[:])
		if user2 != nil {
			p.Users = append(p.Users, user2)
		}
	}
	p.Messages = make(chan *Message)
	return p, nil
}

//serialize the struct to bytes
//func (p *Group) ToBytes() ([]byte, error) {
func (p *Group) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}


	buf.Write(p.Name[:])
	buf.WriteByte(byte('\n'))

	buf.Write(p.Owner.Name[:])

	for _, user := range p.Users {
		buf.Write(user.Name[:])
	}

	return buf.Bytes(), nil
}

func (p *Group) HashKey() interface{} {
	return p.Name
}

func (p *Group) AddUser(user *User) {
	p.Users = append(p.Users, user)
}

func NewGroup(name string, owner *User) *Group {

	if len(name) > 16 {
		return nil
	}

	p := &Group{
		Owner: owner,
		Users: []*User{},
	}
	copy(p.Name[:], []byte(name))
	p.Messages = make(chan *Message)

	return p
}
