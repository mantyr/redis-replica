package redis_replica

import (
	"fmt"
	"bufio"
	"sync"
	"bytes"
)

var (
	// ChannelBufferDefault определяет размер очерди по умолчанию
	ChannelBufferDefault = 100
)

type Clienter interface {
	Do() error
	GetChannel() <-chan *Command
	Write([]byte) error
	WriteCommand(*Command) error
}

type Client struct {
	tr            Transport
	body          *bufio.Reader

	// commandWriterChan канал для передачи комманд на сервер
	commandWriterChan    chan *Command
	// commandReaderChan канал для чтения комманд с сервера
	commandReaderChan    chan *Command

	// Определяем статус клиента
	sync.RWMutex
	runStatus       bool
	closeNotifyChan []chan struct{}
	closeStatus     bool
}

// NewClient создаёт клиент к транспорту
func NewClient(tr Transport, channelBuffer int) (*Client, error) {
	if tr == nil {
		return nil, fmt.Errorf("Expected Transport")
	}
	if channelBuffer < 0 {
		channelBuffer = ChannelBufferDefault
	}
	
	c := new(Client)
	c.tr = tr
	c.commandWriterChan = make(chan *Command, channelBuffer)
	c.commandReaderChan = make(chan *Command, channelBuffer)
	return c, nil
}

// Connect возвращает клиент подключенный к внешнему источнику
func Connect(host string, port int, channelBuffer int) (*Client, error) {
	tr, err := NewNetTransport(host, port)
	if err != nil {
		return nil, fmt.Errorf("Bad Transport, %q", err)
	}
	c, err := NewClient(tr, channelBuffer)
	return c, err
}

// ConnectReader возвращает клиент подключенный к bufio.Reader
func ConnectReader(reader *bufio.Reader, channelBuffer int) (*Client, error) {
	tr, err := NewReaderTransport(reader)
	if err != nil {
		return nil, fmt.Errorf("Bad Transport, %q", err)
	}
	c, err := NewClient(tr, channelBuffer)
	return c, nil
}

// ConnectString возвращает клиент подключенный к строке превращённой в bufio.Reader
func ConnectString(body string, channelBuffer int) (*Client, error) {
	r := bufio.NewReader(bytes.NewBufferString(body))
	return ConnectReader(r, channelBuffer)
}

// Do запускает процесс чтения (запускает только один раз)
func (c *Client) Do() error {
	c.Lock()
	defer c.Unlock()
	
	if c.runStatus {
		return nil
	}
	
	var err error
	c.body, err = c.tr.Dial()
	if err != nil {
		return fmt.Errorf("Error Client.Do, %q", err)
	}

	c.runStatus = true
	go c.writer()
	go c.reader()
	return nil
}

// Client закрывает соединение с сервером если оно есть, прекращает поток данных
func (c *Client) Close() error {
	if !c.runStatus {
		return fmt.Errorf("Client is not running")
	}
	err := c.tr.Close()
	if err != nil {
		return err
	}
	c.setCloseStatus()
	return nil
}

// Write обеспечивает синхронную запись в сервер
func (c *Client) Write(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("Expected data []byte, but actial nil")
	}
	return c.tr.Write(data)
}

// WriteCommand обеспечивает синхронную запись в сервер
func (c *Client) WriteCommand(command *Command) error {
	if command == nil {
		return fmt.Errorf("Expected Command, but actual nil")
	}
	return c.tr.Write(command.GetRawData())
}

// writer читает комманды из канала и передаёт их источнику
func (c *Client) writer() {
	defer c.Close()

	var command *Command
	var err     error
	for command = range c.commandWriterChan {
		err = c.WriteCommand(command)
		if err != nil {
			break
		}
	}
	// завершится когда закроется канал c.commandWriterChan
}

// reader читает комманды из источника и записывает их в канал
func (c *Client) reader() {
	defer c.Close()

	var command *Command
	var err     error
	for {
		// сюда вставить проверку на close
		command, err = readRedisCommand(c.body)
		if err != nil {
			break
		}
		c.commandReaderChan <- command
	}
}

func (c *Client) GetChannel() (<-chan *Command) {
	return c.commandReaderChan
}