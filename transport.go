package redis_replica

import (
	"net"
	"fmt"
	"bufio"
)

var (
	BufSizeDefault int = 16384
)

// Transport описывает интерфейс источника данных
type Transport interface {
	// Dial устанавливает соединение и возвращает bufio.Reader
	Dial() (*bufio.Reader, error)

	// Write передаёт данные источнику
	Write([]byte) error

	// Close закрывает канал
	Close() error
}

// =====================================================

// NetTransport implementation Transport intrface
type NetTransport struct {
	host string
	port int
	conn net.Conn
}

// NewNetTransport возвращает новый транспорт на основе tcp коннекта
func NewNetTransport(host string, port int) (*NetTransport, error) {
	if len(host) == 0 || port == 0 {
		return new(NetTransport), fmt.Errorf("Bad host or port in Transport")
	}
	n := new(NetTransport)
	n.host = host
	n.port = port
	return n, nil
}

// Dial подключается по сети и возвращает bufio.Reader
func (n *NetTransport) Dial() (*bufio.Reader, error) {
	var err error

	n.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", n.host, n.port))
	if err != nil {
		return &bufio.Reader{}, err
	}
	reader := bufio.NewReaderSize(n.conn, BufSizeDefault)
	return reader, nil
}

// Write передаёт данные серверу
func (n *NetTransport) Write(data []byte) error {
	_, err := n.conn.Write(data)
	return err
}

func (n *NetTransport) Close() error {
	if n.conn == nil {
		return fmt.Errorf("No connect")
	}
	return n.conn.Close()
}

// =====================================================

// ReaderTransport implementation Transport intrface
type ReaderTransport struct {
	body *bufio.Reader
}

// NewReaderTransport возвращает новый транспорт на основе bufio.Reader
func NewReaderTransport(reader *bufio.Reader) (*ReaderTransport, error) {
	r := new(ReaderTransport)
	r.body = reader
	return r, nil
}

// Dial возвращает bufio.Reader
func (r *ReaderTransport) Dial() (*bufio.Reader, error) {
	return r.body, nil
}

// Write не передаёт никуда данные так как источник (сервер) является bufio.Reader и записать в него ничего нельзя
func (r *ReaderTransport) Write(data []byte) error {
	return nil
}

func (r *ReaderTransport) Close() error {
	return nil
}