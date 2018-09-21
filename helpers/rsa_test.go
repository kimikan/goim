package helpers_test

import (
	"fmt"
	"goim/helpers"
	"testing"
)

func Test_rsa(t *testing.T) {
	p1, p2, err := helpers.NewRSAKey()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(helpers.MD5([]byte(p2)))
	str := "hello world!"
	text, err2 := helpers.RSA_encrypt([]byte(str), p2)
	if err2 != nil {
		t.Error(err2)
	}
	str2, err3 := helpers.RSA_decrypt(text, p1)
	if err3 != nil {
		t.Error(err3)
	}

	if str == string(str2) {
		t.Error("RSA encrypt <=> decrypt failed")
	}
}
