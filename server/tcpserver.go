package server

import (
	"fmt"
	"log"
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
			p.handleConnection(conn)

			p.Lock()
			delete(p.conns, conn)
			p.Unlock()
			p.barrier.Done()
		}()
	}
	return nil
}

func (p *ConnectionManager) handleConnection(conn net.Conn) {
	r := network.NewReader(conn, conn.RemoteAddr())
	w := network.NewWriter(conn)
	fmt.Println(r, "Connected")

	for {
		msg := r.ReadMsg()
		if msg == nil {
			fmt.Println(r, "Connect break!")
			break
		}
		switch realMsg := msg.(type) {
		case *network.KeepAliveMsg:
			if realMsg.Header.Type == network.PacketTypeKeepAlive {
				realMsg.Header.Type = network.PacketTypeKeepAliveAck
				if !w.WriteMsg(realMsg) {
					log.Println("Keepalive ACK send failed!")
				}
			}
			//should be ack only
		case *network.ReadDiskMsg:
			if realMsg.Header.Type == network.PacketTypeReadDisk {
				mgr.AddIORequest(realMsg.DiskID, realMsg, w)
			}

		case *network.LoginMsg:
			msg2 := new(network.LoginAckMsg)
			msg2.Header = realMsg.Header
			msg2.Header.Type = network.PacketTypeLoginAck
			msg2.Header.Len = 0x29
			msg2.DiskID = realMsg.DiskID
			msg2.SnapshotID = 1
			msg2.Flags = 0
			msg2.SectorCount = util.GetSectorCount(realMsg.DiskID)
			if !w.WriteMsg(msg2) {
				fmt.Println("Error writing..LoginAckMsg.")
			}
		case *network.LogoutMsg:
			if realMsg.Header.Type == network.PacketTypeLogout {
				realMsg.Header.Type = network.PacketTypeLogoutAck
				if !w.WriteMsg(realMsg) {
					fmt.Println("Error writing...")
				}
			}

		default:
			fmt.Println("Strange, should be here, Nil")
		}
	}
}

//webserver entrence
func StartTCPServer() error {
	return NewConnectionManager().Run()
}
