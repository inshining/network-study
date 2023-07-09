package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)
// timeout 시간을 운영체제에 기댈 수 없기 때문에 자체적으로 타임아웃 설정하는 코드 

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error){
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err: "conncetion time out",
				Name: addr,
				Server: "127.0.0.1",
				IsTimeout: true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

// 한번 연결된 ip와 또 다시 접속하려고 시도하면 안됨. -> 타임 아웃이나 다른 에러 내뱉게됨. 
func TestDialTimeout(t *testing.T){
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil{
		c.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok{
		t.Fatal(err)
	}
	if !nErr.Timeout(){
		t.Fatal("error is not a timeout")
	}
}