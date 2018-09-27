package im

import (
	"errors"
	"goim/helpers"
	"io"
)

const (
	MessageType_Login1 = iota
	MessageType_Login2
	MessageType_Login3
	MessageType_Info
)

type MessageBase struct {
	CheckSum uint32
}

type InfoMsg struct {
	Text string
}

type LoginMsg1 struct {
	MessageBase
	PublicKey string
}

type LoginMsg2 struct {
	MessageBase
	//encrypted message with
	//RSA public key
	EncryptedText []byte
}

type LoginMsg3 struct {
	MessageBase
	//DecryptedText with private
	//key
	DecryptedText []byte
}

//user....
type User struct {
}

func WriteInfoMessage(w io.Writer, msg string) error {
	m := &InfoMsg{
		Text: msg,
	}
	bs, err := helpers.Marshal(m)
	if err != nil {
		return err
	}
	l, err2 := helpers.WriteMessage(w, MessageType_Info, bs)
	if err2 != nil {
		return err2
	}
	if l != len(bs) {
		return errors.New("write failed!")
	}

	return nil
}

func GetInfoMessage(r io.Reader) (*InfoMsg, error) {
	t, buf, e := helpers.ReadMessage(r)
	if e != nil {
		return nil, e
	}
	if t != MessageType_Info {
		return nil, errors.New("Wrong message sequence!")
	}
	var m1 InfoMsg
	err := helpers.UnMarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}
