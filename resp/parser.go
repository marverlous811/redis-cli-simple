package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func readline(br *bufio.Reader) ([]byte, error) {
	p, err := br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, fmt.Errorf("bad resp line")
	}
	return p[:i], nil
}

func Parse(br *bufio.Reader) (interface{}, error) {
	line, err := readline(br)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, fmt.Errorf("not have response")
	}

	switch line[0] {
	case '+':
		return string(line[1:]), nil
	case '-':
		return fmt.Errorf(string(line[1:])), nil
	case ':':
		return strconv.Atoi(string(line[1:]))
	case '$':
		length, e := strconv.Atoi(string(line[1:]))
		if length < 0 || e != nil {
			return nil, e
		}
		p := make([]byte, length)
		_, err := io.ReadFull(br, p)
		if err != nil {
			return nil, err
		}

		if nextLine, err := readline(br); err != nil {
			return nil, err
		} else if len(nextLine) != 0 {
			return nil, fmt.Errorf("bad bulk string format")
		}
		return p, nil
	case '*':
		length, e := strconv.Atoi(string(line[1:]))
		if length < 0 || e != nil {
			return nil, e
		}
		r := make([]interface{}, length)
		for i := range r {
			r[i], err = Parse(br)
			if err != nil {
				return nil, err
			}
		}
		return r, nil
	}

	return nil, fmt.Errorf("unexpected resp")
}
