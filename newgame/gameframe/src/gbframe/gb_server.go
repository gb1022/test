package gbframe

import (
	"fmt"
	"net"
	"sync"
)

type Service struct {
	Conn net.Conn
	//	Conn     net.Conn
	Id       string
	Msg      chan string
	State    bool
	TranData *TransportData
	Wg       *sync.WaitGroup
	//	Sig
}

//var NetConnMap = map[string]*Service{}

func (s *Service) ServiceProcess() {
	s.Wg.Add(1)
	for {
		s.TranData.Prosses()
		//		fmt.Println("dddddddddddddddddddddddd, state:", s.TranData.State)
		if s.TranData.State == false {
			s.State = false
			s.ConnClose()
			return
		}
	}
}

func (s *Service) ConnClose() {
	//	s.Conn.Close()

	s.TranData.Conn.Close()
	s.Wg.Done()
	fmt.Println("Service connClose!!!!!!!!!!!!!!!!!!")
}

//func (s *Service) hanldClient() {
//	//	defer func() {
//	//		s.ConnClose(conn)
//	//	}()
//	//	t := CreateTransportData(s.Conn)
//	//	for {
//	//		buf := make([]byte, 2048)
//	//		n, _ := s.Conn.Read(buf)
//	//		fmt.Println("gbframe hanldClient....t.conn.romateip:", t.Conn.RemoteAddr().String())
//	//		s.TranData <- *t
//	//		fmt.Println("gbframe s.TranData <- *t....t.conn.romateip:", t.Conn.RemoteAddr().String())
//	s.Process()
//	//	}
//}

func CreateService(listener *net.TCPListener, id string) (*Service, error) {
	//	fmt.Println("aaaaaaaaaaaaaaaaa")
	conn, err := Connect(listener)
	//	fmt.Println("bbbbbbbbbbbbbbb")
	if err != nil {
		return nil, err
	}
	t := CreateTransportData(conn)
	//	fmt.Println("ccccccccccccccccccc")
	s := &Service{
		Conn: conn,
		//			Conn:  &conn,
		Msg:      make(chan string),
		Id:       id,
		State:    true,
		TranData: t,
		Wg:       &sync.WaitGroup{},
	}
	//		sess := MakeSession(id)
	//		NetConnMap[sess] = s

	//	go s.hanldClient()
	//	for {
	//		s.TranData.Prosses()
	//		fmt.Println("dddddddddddddddddddddddd")
	//	}
	//		//		NetConnMap[conn].RecvMsg = make(chan MsgStruct)
	//		go s.Process()
	//		fmt.Println("-----------")

	return s, nil
}
