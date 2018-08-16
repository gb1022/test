package gbframe

import (
	"fmt"
	"net"
)

//type TcpNetAddr struct {
//	Tcpaddr *net.TCPAddr
//}
//type TCPNetListener struct {
//	Tcplistener net.TCPListener
//}

//func ResolveAddrTcp(pro, addr string) *TcpNetAddr {
//	//	var netAddr *TcpNetAddr
//	var netaddr *net.TCPAddr
//	var err error
//	switch pro {
//	case "tcp", "tcp4", "tcp6":
//		netaddr, err = net.ResolveTCPAddr(pro, addr)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return nil
//		} else {
//			//			fmt.Println("*Addr:", *netaddr)
//			//			fmt.Println("Addr:", netaddr)
//			//			fmt.Println("&Addr:", &netaddr)
//			//			fmt.Println("type:", reflect.TypeOf(netaddr).String())
//			netAddr := TcpNetAddr{
//				Tcpaddr: netaddr,
//			}
//			return &netAddr
//		}
//	default:
//		return nil
//	}
//}

func ListenTcp(pro string, addr string) (*net.TCPListener, error) {
	//	var tcpnetlistener *TCPNetListener
	netaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP(pro, netaddr)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return listener, nil
}

func Connect(listener *net.TCPListener) (net.Conn, error) {
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
