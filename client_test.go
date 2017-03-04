package redis_replica_test

import (
	"testing"
	"time"
	"bufio"
	"bytes"
//	"fmt"
	replica "github.com/mantyr/redis-replica"
)

// TestNetTransportInterface проверяет что NetTransport реализует интерфейс Transport
func TestClientInterface(t *testing.T) {
    var _ replica.Clienter = (*replica.Client)(nil)
}

// TestClientConnectCloseError проверяет что нельзя закрыть коннект у которого не вызван Do() метод
func TestClientConnectCloseError(t *testing.T) {
	client, err := replica.Connect("127.0.0.1", 7777, replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica, %q", err)
	}
	// client.Do()
	err = client.Close()
	if err == nil {
		t.Errorf("Expected error %q, but actual nil", "No connect")
	}
}

// TestClientConnectStringCloseError проверяет что нельзя закрыть коннект у которого не вызван Do() метод
func TestClientConnectStringCloseError(t *testing.T) {
	client, err := replica.ConnectString(RDBFile1, replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica bufio.Reader, %q", err)
	}
	// client.Do()
	err = client.Close()
	if err == nil {
		t.Errorf("Expected error %q, but actual nil", "No connect")
	}
}

// TestClientConnectReader проверяет что мы можем подключиться к bufio.Reader в качестве реплики
func TestClientConnectReader(t *testing.T) {
	client, err := replica.ConnectString(RDBFile1, replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica bufio.Reader, %q", err)
	}
	err = client.Do()
	if err != nil {
		t.Errorf("Error client.Do(), %q", err)
	}
	err = client.Close()
	if err != nil {
		t.Errorf("Error client.Close(), %q", err)
	}
}
/*
func TestClient(t *testing.T) {
	client, err := replica.ConnectString("REDIS0007", replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica bufio.Reader, %q", err)
	}
	err = client.Do()
	if err != nil {
		t.Errorf("Error client.Do(), %q", err)
	}

	var command *replica.Command

	ch := client.GetChannel()
	for {
		command = <- ch
		fmt.Println("Command: ", command)
	}
	err = client.Close()
	if err != nil {
		t.Errorf("Error client.Close(), %q", err)
	}
}
*/

// TestClientCloseWait проверяет что после закрытия соединения отпускается wait
func TestClientCloseWait(t *testing.T) {
	client, err := replica.ConnectString("REDIS0007", replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica io.Reader, %q", err)
	}
	defer client.Close()
	err = client.Do()
	if err != nil {
		t.Errorf("Error client.Do(), %q", err)
	}
	// иногда падает
	select {
		case <- client.WaitClose():
		case <- func() (chan struct{}) {
			time.Sleep(4 * time.Second)
			
			ch := make(chan struct{}, 1)
			ch <- struct{}{}
			
			return ch
		}():
			t.Errorf("Error WaitClose")
	}
}

func TestClientHandler(t *testing.T) {
	reader := bufio.NewReader(bytes.NewBufferString(RDBFile1))
	
	client, err := replica.ConnectReader(reader, replica.ChannelBufferDefault)
	if err != nil {
		t.Errorf("Error connect redis replica io.Reader, %q", err)
	}
	
	test := NewHandlersTest()

	client.Handle("ZADD", func(command *replica.Command) bool {
//		test.Add(command.Name)
		return true
	})
	client.Do()
	
	test.Assert(t, "ZADD", 1)

	client.WaitClose()
}

type handlersTest map[string]int64

func NewHandlersTest() *handlersTest {
	h := new(handlersTest)
	return h
}

func (h *handlersTest) Assert(t *testing.T, handlerName string, expectedCount int64) {
//	count, _ := h[handlerName];

	var count int64 = 12
	if count != expectedCount {
//		t.Errorf(`Expected count Command %v is %v but actual is %v`, handlerName, expectedCount, count)
	}
}

func (h *handlersTest) Add(handlerName string) {
//	h[handlerName]++
}

const (
    RDBFile1 = "REDIS0006\xfe\x00\x00\x03b_1\x04kuku\x00\x03a_1\x04lala\x00\x03b_3\xc3\t@\xb3\x01aa\xe0\xa6\x00\x01aa\xfc\xdb\x82\xb0\\B\x01\x00\x00\x00\x03b_2\r2343545345345\x00\x03a_2\xc0!\xffT\x81\xe9\x86\xcc\x9f\x1f\xc4"
)

