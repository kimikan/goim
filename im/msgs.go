package im

import (
	"errors"
	"goim/helpers"
	"io"

	"github.com/golang/protobuf/proto"
)

const (
	MessageType_Login1 = 1
	MessageType_Login2 = iota
	MessageType_Login3
	MessageType_Info
	MessageType_Heartbeat
	MessageType_UserProfileRequest
	MessageType_UserProfileResponse
	MessageType_UpdateProfileRequest
	MessageType_UpdateProfileResponse
	MessageType_FriendRequest
	MessageType_FriendResponse
	MessageType_ApproveRequest
	MessageType_ApproveResponse
	MessageType_GetAllFriendRequest
	MessageType_GetAllFriendRequestsResponse
	MessageType_Talk
	MessageType_TalkResponse
	MessageType_GetCachedTextRequest
	MessageType_GetCachedTextResponse
	MessageType_Notification
)

func WriteToClientMessage(w io.Writer, msg proto.Message) error {
	if msg == nil {
		return errors.New("invalid parameters")
	}
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	l, err2 := helpers.WriteMessage(w, MessageType_Notification, bs)
	if err2 != nil {
		return err2
	}
	if l != len(bs) {
		return errors.New("write failed")
	}
	return nil
}

func WriteHeartbeatMessage(w io.Writer) error {
	m := &HeartbeatMsg{}
	bs, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	l, err2 := helpers.WriteMessage(w, MessageType_Heartbeat, bs)
	if err2 != nil {
		return err2
	}
	if l != len(bs) {
		return errors.New("write failed!")
	}

	return nil
}

func WriteInfoMessage(w io.Writer, msg string) error {
	m := &InfoMsg{
		Text: msg,
	}
	bs, err := proto.Marshal(m)
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
	err := proto.Unmarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}
