package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"redis-cli/resp"
)

type Conn struct {
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
}

func NewConnection(addr string) (*Conn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn: c,
		bw:   bufio.NewWriter(c),
		br:   bufio.NewReader(c),
	}, nil
}

func (c *Conn) Do(cmd string, args ...interface{}) (interface{}, error) {
	respCmd := resp.BuildCommand(cmd, args...)
	_, err := c.conn.Write(respCmd)
	if err != nil {
		return nil, err
	}
	return resp.Parse(c.br)
}

func printReply(level int, reply interface{}) {
	switch val := reply.(type) {
	case string:
		fmt.Printf("reply: %s\n", reply)
	case []byte:
		fmt.Printf("reply: %s\n", reply)
	case int64:
		fmt.Printf("reply: %d\n", reply)
	case []interface{}:
		for i, v := range val {
			if i != 0 {
				fmt.Printf("%s", strings.Repeat(" ", level*2))
			}

			printReply(level+1, v)
			if i != len(val)-1 {
				fmt.Printf("\n")
			}
		}
	default:
		fmt.Printf("reply: %v\n", reply)
	}
}

func main() {
	addr := "127.0.0.1:6379"
	reader := bufio.NewReader(os.Stdin)
	c, err := NewConnection(addr)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return
	}

	for {
		fmt.Printf("[%s] >> ", addr)
		text, _ := reader.ReadBytes('\n')
		cmdStr := string(text[:len(text)-1])
		if cmdStr == "exit" {
			break
		}
		cmdArr := strings.Split(cmdStr, " ")
		if len(cmdArr) < 0 {
			continue
		}
		strArg := cmdArr[1:]
		args := make([]interface{}, len(strArg))
		for i := range args {
			args[i] = strArg[i]
		}
		reply, err := c.Do(cmdArr[0], args...)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		printReply(0, reply)
	}
}
