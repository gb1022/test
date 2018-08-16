package gbframe

import (
	"fmt"
	"net"
)

type TransportData struct {
	Conn    net.Conn
	InData  chan []byte
	OutData chan []byte
	State   bool
	//	mu      sync.Mutex // guard the following
}

func (t *TransportData) ReadData() {
	buf := make([]byte, 2048)
	n, _ := t.Conn.Read(buf)
	if n == 0 || buf == nil {
		t.State = false
		return
	}
	t.InData <- buf[:n]
	//	fmt.Println("eeeeeeeeeeeeeeee")
}

func (t *TransportData) WriteData() {
	for {
		select {
		case d := <-t.OutData:
			fmt.Println("gbframe writedata:", d)
			t.Conn.Write(d)
		}
	}
}

func (t *TransportData) Prosses() {
	//	for {
	t.ReadData()
	go t.WriteData()

	//	}

}
func CreateTransportData(conn net.Conn) *TransportData {
	//	var outdata []byte
	newTransportData := &TransportData{
		Conn:    conn,
		InData:  make(chan []byte),
		OutData: make(chan []byte),
		State:   true,
	}
	return newTransportData
}
