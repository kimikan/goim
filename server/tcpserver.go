package server

import (
	"errors"
	"fmt"
	"goim/dispatcher"
	"goim/helpers"
	"goim/im"
	"net"
	"sync"
	"time"
)

type ConnectionManager struct {
	sync.RWMutex
	conns map[net.Conn]bool
	//to ensure all of the connections was
	//proper closed before process quit
	barrier sync.WaitGroup
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns: make(map[net.Conn]bool),
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
	err := dispatcher.HandleLogin(conn)
	if err != nil {
		return errors.New("login failed, close connection!")
	}
	err = im.WriteInfoMessage(conn, "Welcome abord")
	if err != nil {
		return err
	}

	for {
		t, buf, e := helpers.ReadMessage(conn)
		if e != nil {
			fmt.Println(e)
			break
		}
		isMgr, e2 := dispatcher.HandleUserMgr(conn, t, buf)
		if e2 != nil {
			fmt.Println(e)
			break
		}
		if isMgr {
			continue
		}

		fmt.Println(t, buf)
		/*
			switch realMsg := msg.(type) {
			default:
				fmt.Println("Strange, should be here, Nil")
			}
		*/
	}

	return nil
}

//webserver entrence
func StartTCPServer() error {
	return NewConnectionManager().Run()
}
