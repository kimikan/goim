package dispatcher

import (
	"time"
)

//the structures defined in this file
//used to communicate with client via messagebus

//while a recieve fromuser's friend request
//do this, if online
//topic is myself's userid
type NotificationFriendRequest struct {
	FromUserID string
	HelloMsg   string
}

//topic is other side's userid
type NotificationApprove struct {
	ApprovedFriendID string
}

type NotificationText struct {
	FromUserID  string
	Content     string
	ArrivedTime time.Time
}
