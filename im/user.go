package im

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
)

//definition user
type User struct {
	ID          uuid.UUID
	Name        string
	Password    string
	Description string
}

//deserialize the buf to struct
func ParseUser(bs []byte) (*User, error) {
	p := &User{}
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

	p.Password, err = buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	p.Password = strings.TrimRight(p.Password, "\n")
	p.Description, err = buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	p.Description = strings.TrimRight(p.Description, "\n")

	return p, nil
}

//serialize the struct to bytes
func (p *User) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	b, err := p.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf.Write(b)
	buf.WriteString(p.Name)
	buf.WriteByte(byte('\n'))
	buf.WriteString(p.Password)
	buf.WriteByte(byte('\n'))
	buf.WriteString(p.Description)
	buf.WriteByte(byte('\n'))

	return buf.Bytes(), nil
}

//implements the hash function
func (p *User) HashKey() interface{} {
	return p.ID
}

func (p *User) SetDescription(strs string) {
	p.Description = strs
}

func NewUser(name string, pass string) *User {
	id, err := uuid.NewUUID()

	if err != nil {
		return nil
	}

	p := &User{
		ID:          id,
		Name:        name,
		Password:    pass,
		Description: "",
	}

	return p
}
