package gbframe

import (
	"crypto/md5"
	"encoding/hex"

	// "fmt"
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
	n, err := t.Conn.Read(buf)
	if n == 0 || buf == nil || err != nil {
		Logger_Error.Println("ReadData err:", err)
		// t.State <- false
		t.State = false
		return
	}
	t.InData <- buf[:n]
	// Logger_Info.Println("eeeeeeeeeeeeeeee,n:", n, " buf:", buf[:n])
}

func (t *TransportData) WriteData() {
	for {
		select {
		case d := <-t.OutData:
			// fmt.Println("gbframe writedata:", d)
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
		InData:  make(chan []byte, 100),
		OutData: make(chan []byte, 100),
		State:   true,
	}
	return newTransportData
}

func MakeSession(str string, pass string) string {
	pass_byte := []byte(pass)
	if pass == "" {
		pass_byte = nil
	}
	h := md5.New()
	h.Write([]byte(str))
	cipher := h.Sum(pass_byte)
	md5_str := hex.EncodeToString(cipher)
	return md5_str
}
