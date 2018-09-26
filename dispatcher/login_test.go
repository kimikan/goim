package dispatcher_test

import (
	"errors"
	"fmt"
	"goim/helpers"
	"goim/im"
	"net"
	"testing"
)

func writeLogin1(conn net.Conn, pub string) error {
	m1 := im.LoginMsg1{
		PublicKey: pub,
	}
	bs1, ex := im.Marshal(&m1)
	if ex != nil {
		return ex
	}
	n, err3 := helpers.WriteMessage(conn, im.MessageType_Login1, bs1)
	if n != len(bs1) {
		return errors.New("write failed!")
	}
	return err3
}

func writeLogin3(conn net.Conn, private string, encryptedText []byte) error {
	raw, err2 := helpers.RSA_decrypt(encryptedText, private)
	if err2 != nil {
		return err2
	}
	m3 := im.LoginMsg3{
		DecryptedText: raw,
	}
	bs3, ex := im.Marshal(&m3)
	if ex != nil {
		return ex
	}
	n, err3 := helpers.WriteMessage(conn, im.MessageType_Login3, bs3)
	if n != len(bs3) {
		return errors.New("write failed!")
	}
	return err3
}

//verify with source
func getLogin2(conn net.Conn) (*im.LoginMsg2, error) {
	t, buf, e := helpers.ReadMessage(conn)
	if e != nil {
		return nil, e
	}
	if t != im.MessageType_Login2 {
		return nil, errors.New("Wrong message sequence!")
	}
	var m1 im.LoginMsg2
	err := im.UnMarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}

func Test_login(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Error(err)
		return
	}
	if conn != nil {
		defer conn.Close()
	}
	pri, pub, e := helpers.NewRSAKey()
	if e != nil {
		t.Error(e)
		return
	}
	err = writeLogin1(conn, pub)
	if err != nil {
		t.Error(err)
		return
	}
	m2, e2 := getLogin2(conn)
	if e2 != nil {
		t.Error(e2)
		return
	}
	e3 := writeLogin3(conn, pri, m2.EncryptedText)
	if e3 != nil {
		t.Error(e3)
		return
	}
	info, e4 := im.GetInfoMessage(conn)
	if e4 != nil {
		t.Error(e4)
		return
	}
	fmt.Println(info)
	//t.Error(info)
}
