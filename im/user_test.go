package im_test

import (
	"fmt"
	"goim/helpers"
	"goim/im"
	"testing"
)

func Test_toid(t *testing.T) {
	_, pub, err := helpers.NewRSAKey()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pub)
	id1 := im.KeyToUserID(pub)
	id2 := im.KeyToUserID("----\n" + pub)
	if id1 != id2 {
		t.Error("fucking")
	}
}
