package dispatcher

import (
	"bytes"
	"errors"
	"fmt"
	"goim/helpers"
	"goim/im"
	"math/rand"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

func randBytes() []byte {
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		result = append(result, byte(r.Intn(255)))
	}
	return result
}

func getLogin1(conn net.Conn) (*im.LoginMsg1, error) {
	t, buf, e := helpers.ReadMessage(conn)
	if e != nil {
		return nil, e
	}
	if t != im.MessageType_Login1 {
		return nil, errors.New("Wrong message sequence!")
	}
	var m1 im.LoginMsg1
	err := proto.Unmarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}

func replyLogin2(conn net.Conn, key string, bys []byte) error {
	text, err2 := helpers.RSA_encrypt(bys, key)
	if err2 != nil {
		return err2
	}
	m2 := im.LoginMsg2{
		EncryptedText: text,
	}
	bs3, ex := proto.Marshal(&m2)
	if ex != nil {
		return ex
	}
	n, err3 := helpers.WriteMessage(conn, im.MessageType_Login2, bs3)
	if n != len(bs3) {
		return errors.New("write failed!")
	}
	return err3
}

//verify with source
func getLogin3(conn net.Conn) (*im.LoginMsg3, error) {
	t, buf, e := helpers.ReadMessage(conn)
	if e != nil {
		return nil, e
	}
	if t != im.MessageType_Login3 {
		return nil, errors.New("Wrong message sequence!")
	}
	var m1 im.LoginMsg3
	err := proto.Unmarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}

//return userid, error
func HandleLogin(conn net.Conn) (string, error) {
	m1, e := getLogin1(conn)
	if e != nil {
		return "", e
	}
	//if not register
	//should submit regitstion information
	//via webserver
	_, err := im.GetUserInfoByPubKey(m1.PublicKey)
	if err != nil {
		return "", err
	}
	bys := randBytes()
	e2 := replyLogin2(conn, m1.PublicKey, bys)
	if e2 != nil {
		return "", e2
	}
	m3, e3 := getLogin3(conn)
	if e3 != nil {
		return "", e3
	}
	if bytes.Compare(m3.DecryptedText, bys) != 0 {
		return "", errors.New("Wrong password!")
	}
	fmt.Println("Success")
	id := im.KeyToUserID(m1.PublicKey)
	return id, nil
}
