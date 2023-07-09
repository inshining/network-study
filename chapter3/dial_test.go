package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T){
	listener, err := net.Listen("tcp", "127.0.0.1:0");
	if (err != nil){
		t.Fatal(err);
	}

	// 고루틴에 필요한 채널 생성 
	done := make(chan struct{})

	go func(){
		// 연결 끊기 전에 채널 초기화
		defer func() { done <- struct{}{}}()

		for {
			// 리스너 연결 성립 
			// 고 언어에서 핸드 쉐이킹 사전에 해줌 
			conn, err := listener.Accept();
			if err != nil{
				t.Log(err)
				return
			}

			go func(c net.Conn){
				// 커넥션 후 핸들링에 관한 고루틴 
				defer func(){
					c.Close();
					done <- struct{}{};
				}()

				buf := make([]byte, 1024)

				for {
					n, err := c.Read(buf)
					if (err != nil){
						if (err != io.EOF){
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil{
		t.Fatal(err)
	}

	conn.Close()
	<- done
	listener.Close()
	<- done
}