package im

import (
	"bytes"
	"encoding/gob"
	"log"
)

const (
	MessageType_Login1 = iota
	MessageType_Login2
	MessageType_Login3
)

type MessageBase struct {
	CheckSum uint32
}

type LoginMsg1 struct {
	MessageBase
	PublicKey string
}

func UnMarshal(b []byte, msg interface{}) error {
	var buf = bytes.Buffer{}
	buf.Write(b)
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(msg)
	if err != nil {
		log.Fatal("decode:", err)
		return err
	}
	return nil
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

func Marshal(o interface{}) ([]byte, error) {
	var buf bytes.Buffer
	// Create an encoder and send a value.
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	if err != nil {
		log.Fatal("encode:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
