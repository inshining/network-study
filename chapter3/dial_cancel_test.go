package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// 연결이 어려운 노드에 일부러 연결 시도하면 데드라인에 의해서 연결 시도가 종료될 수 있으나
// 아래 코드에서 일정 시간 (1초) 지나기 전에 연결 시도가 안되면 에러나면 cancel 함수 호출
// 컨텍스트 이용하여서 cancel 함수 생성한다. 
// TODO: 여기서 컨텍스트의 의미와 설명이 완전히 이해하지 못함. 
func TestDialContextCancel(t *testing.T){
	ctx, cancel := context.WithCancel(context.Background())
	sync := make(chan struct{})

	go func(){
		defer func() {sync <- struct{}{}}()

		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil{
			t.Log(err)
			return 
		}

		conn.Close()
		t.Error("connection did not time out")
	}()

	cancel()
	<-sync

	if ctx.Err() != context.Canceled{
		t.Errorf("expected canceled context; actual; %q", ctx.Err())
	}
}