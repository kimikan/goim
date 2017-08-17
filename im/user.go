package im

import (
	"bytes"
	"strings"
	"time"
)

/*
 * Message details
 * to share the data among
 * different users & groups
 */
type Message struct {
	Time    time.Time
	Group   *Group
	From    *User
	Content string
}

//definition user
type User struct {
	Name        [16]byte
	Password    string
	Description string

	Messages chan *Message
}

//start a discussion session
//againest a user
func (p *User) TalkToUser(user *User, msg string) {
	message := &Message{
		Time:    time.Now(),
		Group:   nil,
		From:    p,
		Content: msg,
	}

	if user != nil {
		user.Messages <- message
	}
}

//submit a post in a group discusstion
func (p *User) TalkInGroup(group *Group, msg string) {
	message := &Message{
		Time:    time.Now(),
		Group:   group,
		From:    p,
		Content: msg,
	}

	if group != nil {
		group.Messages <- message

		if len(group.Messages) > 100 {
			<-group.Messages
		}

		for _, user := range group.Users {
			user.Messages <- message
		}
	}
}

//deserialize the buf to struct
func ParseUser(bs []byte) (*User, error) {
	p := &User{}
	buf := bytes.NewBuffer(bs)

	len, err := buf.Read(p.Name[:])
	if len != 16 || err != nil {
		return nil, err
	}
	/*
		p.Name, err = buf.ReadString('\n')
		if err != nil {
			return nil, err
		}
		p.Name = strings.TrimRight(p.Name, "\n") */
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
	p.Messages = make(chan *Message)

	return p, nil
}

//serialize the struct to bytes
func (p *User) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}

	buf.Write(p.Name[:])
	buf.WriteByte(byte('\n'))
	buf.WriteString(p.Password)
	buf.WriteByte(byte('\n'))
	buf.WriteString(p.Description)
	buf.WriteByte(byte('\n'))

	return buf.Bytes(), nil
}

//implements the hash function
func (p *User) HashKey() interface{} {
	return p.Name[:]
}

func (p *User) SetDescription(strs string) {
	p.Description = strs
}

func NewUser(name string, pass string) *User {

	if len(name) > 16 {
		return nil
	}

	p := &User{
		Password:    pass,
		Description: "",
	}
	copy(p.Name[:], []byte(name))
	p.Messages = make(chan *Message)

	return p
}
