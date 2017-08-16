package im

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
)

type Group struct {
	ID   uuid.UUID
	Name string

	//for privilege approve etc
	Owner *User

	Users []*User
}

//deserialize the buf to struct
func ParseGroup(store Store, bs []byte) (*Group, error) {
	p := NewGroup("", nil)
	buf := bytes.NewBuffer(bs)

	tmp := [16]byte{}
	len, err := buf.Read(tmp[:])
	if len != 16 || err != nil {
		return nil, err
	}
	err = p.ID.UnmarshalBinary(tmp[:])
	if err != nil {
		return nil, err
	}

	p.Name, err = buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	p.Name = strings.TrimRight(p.Name, "\n")
	//fmt.Println([]byte(p.Name))
	len, err = buf.Read(tmp[:])
	if len != 16 || err != nil {
		return nil, err
	}
	uid := uuid.UUID{}
	err = uid.UnmarshalBinary(tmp[:])
	if err != nil {
		return nil, err
	}
	p.Owner = store.GetUser(&uid)

	for {
		len, err = buf.Read(tmp[:])
		if len != 16 || err != nil {
			break
		}
		uid := uuid.UUID{}
		err = uid.UnmarshalBinary(tmp[:])
		if err != nil {
			break
		}

		user2 := store.GetUser(&uid)
		if user2 != nil {
			p.Users = append(p.Users, user2)
		}
	}
	return p, nil
}

//serialize the struct to bytes
//func (p *Group) ToBytes() ([]byte, error) {
func (p *Group) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	b, err := p.ID.MarshalBinary()

	if err != nil {
		return nil, err
	}

	buf.Write(b)
	buf.WriteString(p.Name)
	buf.WriteByte(byte('\n'))
	b, err = p.Owner.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	for _, user := range p.Users {
		b, err = user.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	return buf.Bytes(), nil
}

func (p *Group) HashKey() interface{} {
	return p.ID
}

func (p *Group) AddUser(user *User) {
	p.Users = append(p.Users, user)
}

func NewGroup(name string, owner *User) *Group {
	id, err := uuid.NewUUID()

	if err != nil {
		return nil
	}

	p := &Group{
		ID:    id,
		Name:  name,
		Owner: owner,
		Users: []*User{},
	}

	return p
}
