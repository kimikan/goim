package dispatcher

import (
	"goim/helpers"
	"goim/im"
	"net"

	"github.com/golang/protobuf/proto"
)

func handleUserProfileRequest(conn net.Conn, buf []byte, userid string) error {
	var m im.UserProfileRequestMsg
	e := proto.Unmarshal(buf, &m)
	if e != nil {
		return e
	}
	friend := m.GetUserID()
	info, e2 := im.GetUserInfoByID(friend)
	if e2 != nil {
		return e2
	}

	response := &im.UserProfileResponseMsg{}
	response.Result = im.UserProfileResponseMsg_Success
	response.UserID = friend
	response.PublicKey = info.Key
	response.DisplayName = info.DisplayName
	response.Description = info.Description
	response.Avatar = info.Avatar

	if m.GetUserID() != userid {
		contact, e3 := info.GetContactByID(userid)
		if e3 != nil {
			return e3
		}
		if contact == nil {
			//means have no permission to get the specific profile
			response.Reset()
			response.Result = im.UserProfileResponseMsg_NoPermission
		}
	}

	bs, err := proto.Marshal(response)
	if err != nil {
		return err
	}
	_, ex := helpers.WriteMessage(conn, im.MessageType_UserProfileResponse, bs)
	return ex
}

func writeMessage(conn net.Conn, msgType uint32, msg proto.Message) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, ex := helpers.WriteMessage(conn, msgType, bs)
	return ex
}

func handleUpdateProfileRequest(conn net.Conn, buf []byte, userid string) error {
	var m im.UpdateProfileRequestMsg
	e := proto.Unmarshal(buf, &m)
	if e != nil {
		return e
	}
	info, e2 := im.GetUserInfoByID(userid)
	if e2 != nil {
		response := &im.UpdateProfileResponseMsg{
			Result: im.UpdateProfileResponseMsg_UserNotExists,
		}
		return writeMessage(conn, im.MessageType_UpdateProfileResponse, response)
	}

	info.Avatar = m.Avatar
	info.DisplayName = m.DisplayName
	info.Description = m.Description
	e3 := im.SetUserInfo(info)
	response := &im.UpdateProfileResponseMsg{
		Result: im.UpdateProfileResponseMsg_Success,
	}
	if e3 != nil {
		response.Result = im.UpdateProfileResponseMsg_Unspecified
	}
	return writeMessage(conn, im.MessageType_UpdateProfileResponse, response)
}

//parse request msg
//get from's info
//if not exists, return not exists
//info.setrequest(),  done
func handleFriendRequest(conn net.Conn, buf []byte, userid string) error {
	var m im.FriendRequestMsg
	e := proto.Unmarshal(buf, &m)
	if e != nil {
		return e
	}
	info, e2 := im.GetUserInfoByID(m.FriendUserID)
	if e2 != nil {
		response := &im.FriendResponseMsg{
			Result: im.FriendResponseMsg_UserNotExists,
		}
		return writeMessage(conn, im.MessageType_FriendResponse, response)
	}
	e3 := info.AddNewRequestByID(m.FriendUserID, m.HelloMsg)
	if e3 != nil {
		return e3
	}
	response := &im.FriendResponseMsg{
		Result: im.FriendResponseMsg_Success,
	}
	return writeMessage(conn, im.MessageType_FriendResponse, response)
}

//update self's request's state
//add new contact
//if need to judge
func handleApproveRequest(conn net.Conn, buf []byte, userid string) error {
	var m im.ApproveRequestMsg
	e := proto.Unmarshal(buf, &m)
	if e != nil {
		return e
	}
	info, e2 := im.GetUserInfoByID(m.FriendUserID)
	if e2 != nil {
		response := &im.ApproveResponseMsg{
			Result: im.ApproveResponseMsg_RequestNotExists,
		}
		return writeMessage(conn, im.MessageType_ApproveResponse, response)
	}
	myinfo, e3 := im.GetUserInfoByID(userid)
	if e3 != nil {
		return e3
	}
	req, e4 := myinfo.GetRequestByID(m.FriendUserID)
	if e4 != nil {
		response := &im.ApproveResponseMsg{
			Result: im.ApproveResponseMsg_RequestNotExists,
		}
		return writeMessage(conn, im.MessageType_ApproveResponse, response)
	}
	if req.IsApproved {
		response := &im.ApproveResponseMsg{
			Result: im.ApproveResponseMsg_AlreadyApproved,
		}
		return writeMessage(conn, im.MessageType_ApproveResponse, response)
	}
	e5 := myinfo.ApproveRequestByID(m.FriendUserID)
	if e5 != nil {
		return e5
	}
	e6 := myinfo.SetContact(info.Key)
	if e6 != nil {
		return e6
	}
	response := &im.ApproveResponseMsg{
		Result: im.ApproveResponseMsg_Success,
	}
	return writeMessage(conn, im.MessageType_ApproveResponse, response)
}

//handleGetAllRequests
func handleGetAllRequests(conn net.Conn, buf []byte, userid string) error {
	info, e2 := im.GetUserInfoByID(userid)
	if e2 != nil {
		return e2
	}
	rs, e3 := info.GetAllRequests()
	if e3 != nil {
		response := &im.GetAllRequestsResponseMsg{
			Result: im.GetAllRequestsResponseMsg_Unspecified,
		}
		return writeMessage(conn, im.MessageType_GetAllFriendRequestsResponse, response)
	}

	response := &im.GetAllRequestsResponseMsg{
		Result: im.GetAllRequestsResponseMsg_Success,
	}
	response.Requests = []*im.Request{}
	for _, v := range rs {
		response.Requests = append(response.Requests, &im.Request{
			FromID:     v.FromUserID,
			HelloMsg:   v.HelloMsg,
			IsApproved: v.IsApproved,
			Time:       v.ModifiedTime.String(),
		})
	}
	return writeMessage(conn, im.MessageType_GetAllFriendRequestsResponse, response)
}

//request userprofile
//update self's profile
//shenqing haoyou
//approve friends
//get all friend requests
//everyone should have a reply
func HandleUserMgr(bus helpers.MessageBus, conn net.Conn, t uint32, buf []byte, userid string) (bool, error) {
	switch t {
	case im.MessageType_ApproveRequest:
		return true, handleApproveRequest(conn, buf, userid)
	case im.MessageType_FriendRequest:
		return true, handleFriendRequest(conn, buf, userid)
	case im.MessageType_UpdateProfileRequest:
		return true, handleUpdateProfileRequest(conn, buf, userid)
	case im.MessageType_UserProfileRequest:
		return true, handleUserProfileRequest(conn, buf, userid)
	case im.MessageType_GetAllFriendRequest:
		return true, handleGetAllRequests(conn, buf, userid)
	}
	return false, nil
}
