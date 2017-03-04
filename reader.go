package redis_replica

import (
	"io"
	"fmt"
	"strings"
	"strconv"
	"bufio"
)

type CommandType int

const (
	CommandEmpty CommandType = iota   // Пустая комманда
	CommandReply                      // Ответ сервера
	CommandList                       // Блок с коммандами
	CommandBulk                       // Собственно сама комманда
)

type Command struct {
    Type     CommandType    // тип комманды
    Reply    string         // ответ сервера
    Key      string         // ключ
    Args     []string       // аргументы (по сути все параметры после ключа)
    Commands []*Command      // подкомманды если тип CommandList
}

func NewCommand() *Command {
	c := new(Command)
	return c
}

// GetRawData возвращает бинарный формат команды
func (c *Command) GetRawData() []byte {
	return []byte("")
}

// Parse парсит io.Reader на предмет комманды
func (c *Command) Parse(reader bufio.Reader) error {
	header, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if header == "\n" || header == "\r\n" {
		c.Type = CommandEmpty
		return nil
	}

	if strings.HasPrefix(header, "+") {
		c.Type  = CommandReply
		c.Reply = strings.TrimSpace(header[1:])
		return nil
	}
	if strings.HasPrefix(header, "$") {
		bulkSize, err := strconv.ParseInt(strings.TrimSpace(header[1:]), 10, 64)
		if err != nil {
			return fmt.Errorf("Unable to decode bulk size: %v", err)
		}
		c.Type     = CommandBulk
//		c.BulkSize = bulkSize
		_ = bulkSize
		// тут идёт rdb.FilterRDB()
		return nil
	}
	if strings.HasPrefix(header, "*") {
		commandCount, err := strconv.Atoi(strings.TrimSpace(header[1:]))
		if err != nil {
			return fmt.Errorf("Unable to parse command length: %v", err)
		}
		c.Type     = CommandList
		c.Commands = make([]*Command, commandCount)
		
		// Дочитываем комманды
		for i := range c.Commands {
			command := NewCommand()
			err = command.ParseSubCommand(reader)
			if err != nil {
				return fmt.Errorf("Bad command in Command list, %v", err)
			}
			c.Commands[i] = command
		}
		return nil
	}
	return fmt.Errorf("Bad command")
}

// ParseSubCommand возвращает подкомманду группы
func (c *Command) ParseSubCommand(reader bufio.Reader) error {
	c.Type = CommandBulk

	header, err := reader.ReadString('\n')
	if !strings.HasPrefix(header, "$") || err != nil {
		return fmt.Errorf("Failed to read command: %v", err)
	}

	// длина отведённая под аргументы
	argSize, err := strconv.Atoi(strings.TrimSpace(header[1:]))
	if err != nil {
		return fmt.Errorf("Unable to parse argument length: %v", err)
	}
	argument := make([]byte, argSize)
	_, err = io.ReadFull(&reader, argument)
	if err != nil {
		return fmt.Errorf("Failed to read argument: %v", err)
	}
	c.Args = strings.Split(string(argument), " ")

	// дополнительные данные аргументов
	header, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Failed to read argument: %v", err)
	}
	c.Args = append(c.Args, strings.Split(header, " ")...)
	return nil
}

// readRedisCommand читает формат бинарный поток и возвращает сформированную комманду
func readRedisCommand(r io.Reader) (c *Command, err error) {
/*
	if strings.HasPrefix(header, "*") {
		cmdSize, err := strconv.Atoi(strings.TrimSpace(header[1:]))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse command length: %v", err)
		}

		result := &redisCommand{raw: []byte(header), command: make([]string, cmdSize)}

		for i := range result.command {
			header, err = reader.ReadString('\n')
			if !strings.HasPrefix(header, "$") || err != nil {
				return nil, fmt.Errorf("Failed to read command: %v", err)
			}

			result.raw = append(result.raw, []byte(header)...)

			argSize, err := strconv.Atoi(strings.TrimSpace(header[1:]))
			if err != nil {
				return nil, fmt.Errorf("Unable to parse argument length: %v", err)
			}

			argument := make([]byte, argSize)
			_, err = io.ReadFull(reader, argument)
			if err != nil {
				return nil, fmt.Errorf("Failed to read argument: %v", err)
			}

			result.raw = append(result.raw, argument...)

			header, err = reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("Failed to read argument: %v", err)
			}

			result.raw = append(result.raw, []byte(header)...)

			result.command[i] = string(argument)
		}

		return result, nil
	}

	return &redisCommand{raw: []byte(header), command: []string{strings.TrimSpace(header)}}, nil
*/
	return nil, fmt.Errorf("EOF")
}