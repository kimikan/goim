package dispatcher

import (
	"errors"
	"goim/db"
	"goim/helpers"
	"goim/im"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

func DispatchTextMsg(bus helpers.MessageBus, conn net.Conn, m *im.Text, userid string) error {
	if bus == nil || conn == nil || m == nil {
		return errors.New("invalid parameters")
	}
	info, err := im.GetUserInfoByID(m.UserId)
	if err != nil {
		response := &im.TextMsgResponse{
			Result: im.TextMsgResponse_UserNotExists,
		}
		return writeMessage(conn, im.MessageType_TalkResponse, response)
	}
	_, err2 := info.GetContactByID(userid)
	if err2 != nil {
		response := &im.TextMsgResponse{
			Result: im.TextMsgResponse_NotFriend,
		}
		return writeMessage(conn, im.MessageType_TalkResponse, response)
	}
	//if has handle, means, it's online, just forward
	if bus.HasHandle(m.UserId) {
		bus.Publish(m.UserId, &NotificationText{
			FromUserID:  userid,
			Content:     m.Content,
			ArrivedTime: time.Now(),
		})
		return nil
	}

	//otherwise fucking saving in the cache
	err = db.AddMessage(m.UserId, &db.TextItem{})
	if err != nil {
		response := &im.TextMsgResponse{
			Result: im.TextMsgResponse_Unspecified,
		}
		return writeMessage(conn, im.MessageType_TalkResponse, response)
	}
	//decided not to response
	//in order to reduce the message communications
	return nil
}

//text message
func handleTextMessage(bus helpers.MessageBus, conn net.Conn, buf []byte, userid string) error {
	var m im.TextMsg
	e := proto.Unmarshal(buf, &m)
	if e != nil {
		return e
	}
	return DispatchTextMsg(bus, conn, m.Text, userid)
}

func handleGetCache(bus helpers.MessageBus, conn net.Conn, buf []byte, userid string) error {
	items, e := db.GetAllMsgs(userid)
	if e != nil {
		response := &im.GetAllCachedMessagesResponse{
			Result: im.GetAllCachedMessagesResponse_Unspecified,
		}
		return writeMessage(conn, im.MessageType_GetCachedTextResponse, response)
	}
	response := &im.GetAllCachedMessagesResponse{
		Result: im.GetAllCachedMessagesResponse_Success,
	}
	response.Texts = []*im.Text{}
	for _, v := range items {
		response.Texts = append(response.Texts, &im.Text{
			UserId:  v.FromUserID,
			Content: v.Content,
		})
	}
	e2 := writeMessage(conn, im.MessageType_GetCachedTextResponse, response)
	if e2 != nil {
		return e2
	}
	//get once clear it immidiately
	return db.RemoveAllMsgs(userid)
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
