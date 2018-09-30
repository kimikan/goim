package server

import (
	"errors"
	"fmt"
	"goim/db"
	"goim/dispatcher"
	"goim/helpers"
	"goim/im"
	"net"
	"runtime"
	"sync"
	"time"
)

type ConnectionManager struct {
	sync.RWMutex
	conns map[net.Conn]bool
	//to ensure all of the connections was
	//proper closed before process quit
	barrier    sync.WaitGroup
	messageBus helpers.MessageBus
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns:      make(map[net.Conn]bool),
		messageBus: helpers.NewMessageBus(runtime.NumCPU()),
	}
}

func (p *ConnectionManager) Close() {
	//do some clean work
	p.Lock()
	for conn := range p.conns {
		conn.Close()
	}
	p.conns = nil
	p.Unlock()
	p.barrier.Wait()
}

func (p *ConnectionManager) Run() error {
	defer p.Close()

	l, e := net.Listen("tcp", ":9999")
	if e != nil {
		return e
	}
	defer l.Close()
	for {
		conn, err := l.Accept()

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Printf("accept error: %v; retrying in 1s\n", err)
				time.Sleep(time.Second)
				continue
			}
			fmt.Println(err, "Accept")
			break
		}

		p.Lock()
		p.conns[conn] = true
		p.Unlock()
		//barrier
		p.barrier.Add(1)

		go func() {
			defer conn.Close()

			fmt.Println("Ok, New user connected!")
			err := p.handleConnection(conn)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Connection closed!")
			p.Lock()
			delete(p.conns, conn)
			p.Unlock()
			p.barrier.Done()
		}()
	}
	return nil
}

func (p *ConnectionManager) handleConnection(conn net.Conn) error {
	//handle login mechanism
	userid, err := dispatcher.HandleLogin(conn)
	if err != nil {
		return errors.New("login failed, close connection!")
	}
	err = im.WriteInfoMessage(conn, "Welcome abord")
	if err != nil {
		return err
	}

	fn := func(userid string, msg interface{}) {
		switch mx := msg.(type) {
		case *dispatcher.NotificationApprove:
			im.WriteToClientMessage(conn, &im.NotificationApproveMsg{
				UserID: mx.ApprovedFriendID,
			})
		case *dispatcher.NotificationFriendRequest:
			im.WriteToClientMessage(conn, &im.NotificationFriendRequestMsg{
				HelloMsg: mx.HelloMsg,
				UserID:   mx.FromUserID,
			})
		case *dispatcher.NotificationText:
			im.WriteToClientMessage(conn, &im.TextMsg{
				Text: &im.Text{
					Content:  mx.Content,
					UserId:   mx.FromUserID,
					SendTime: mx.ArrivedTime.String(),
				},
			})
		}
	}
	err = p.messageBus.Subscribe(userid, fn)
	if err != nil {
		return err
	}
	defer p.messageBus.Unsubscribe(userid, fn)
	//delivery all of the cached messages
	texts, err2 := db.GetAllMsgs(userid)
	if err2 != nil {
		return err2
	}
	for _, t := range texts {
		ex := im.WriteToClientMessage(conn, &im.TextMsg{
			Text: &im.Text{
				Content:  t.Content,
				UserId:   t.FromUserID,
				SendTime: t.ModifiedTime.String(),
			},
		})
		if ex != nil {
			fmt.Println("Message not delived: ", t)
		}
	}
	err3 := db.RemoveAllMsgs(userid)
	if err3 != nil {
		return err3
	}

	for {
		t, buf, e := helpers.ReadMessage(conn)
		if e != nil {
			fmt.Println(e)
			break
		}
		isMgr, e2 := dispatcher.HandleUserMgr(p.messageBus, conn, t, buf, userid)
		if e2 != nil {
			fmt.Println(e)
			break
		}
		if isMgr {
			continue
		}
		isMsg, e3 := dispatcher.HandleTalkMessage(p.messageBus, conn, t, buf, userid)
		if e3 != nil {
			fmt.Println(e3)
			break
		}
		if isMsg {
			continue
		}

		//others just print
		fmt.Println(t, buf)
	}

	return nil
}

//webserver entrence
func StartTCPServer() error {
	return NewConnectionManager().Run()
}
