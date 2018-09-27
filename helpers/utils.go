package helpers

import (
	"bytes"
	"encoding/gob"
	"log"
)

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
