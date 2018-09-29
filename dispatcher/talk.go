package dispatcher

import (
	"goim/helpers"
	"goim/im"
	"net"
)

func handleTextMessage(bus helpers.MessageBus, conn net.Conn, buf []byte, userid string) error {
	//todo:
	return nil
}

func handleGetCache(bus helpers.MessageBus, conn net.Conn, buf []byte, userid string) error {
	//todo:
	return nil
}

//everyone should have a reply
func HandleTalkMessage(bus helpers.MessageBus, conn net.Conn, t uint32, buf []byte, userid string) (bool, error) {
	if t == im.MessageType_Talk {
		return true, handleTextMessage(bus, conn, buf, userid)
	} else if t == im.MessageType_GetCachedTextRequest {
		return true, handleGetCache(bus, conn, buf, userid)
	}
	return false, nil
}
