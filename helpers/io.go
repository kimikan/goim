package helpers

import (
	"encoding/binary"
	"errors"
	"io"
)

func getUint32(r io.Reader) (uint32, error) {
	var buf [4]byte
	n, err := r.Read(buf[:])
	if n != 4 {
		return 0, errors.New("invalid stream!")
	}
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

//get a new bytes
func getBytes(r io.Reader, num int) ([]byte, error) {
	if num <= 0 {
		return nil, errors.New("invalid parameters!")
	}
	buf := make([]byte, num, num)
	n, err := r.Read(buf)
	if n != num {
		return nil, errors.New("invalid stream")
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//message format[type, len, [len]byte]
func ReadMessage(r io.Reader) (uint32, []byte, error) {
	if r == nil {
		return 0, nil, errors.New("invalid parameter!")
	}
	t, e := getUint32(r)
	if e != nil {
		return 0, nil, e
	}

	l, err := getUint32(r)
	if err != nil {
		return 0, nil, err
	}
	buf, err2 := getBytes(r, int(l))
	if err2 != nil {
		return 0, nil, err2
	}
	return t, buf, nil
}

func writeUint32(w io.Writer, v uint32) error {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	n, err := w.Write(buf[:])
	if n != 4 {
		return errors.New("invalid stream")
	}
	return err
}

//write a messge[type, len, [len]byte]
func WriteMessage(w io.Writer, t uint32, buf []byte) (int, error) {
	if w == nil || buf == nil {
		return 0, errors.New("invalid parameter!")
	}
	len := len(buf)
	if len == 0 {
		return 0, nil
	}
	e := writeUint32(w, t)
	if e != nil {
		return 0, e
	}
	err := writeUint32(w, uint32(len))
	if err != nil {
		return 0, err
	}

	return w.Write(buf)
}
