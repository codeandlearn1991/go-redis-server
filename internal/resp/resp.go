package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type DataType byte

const (
	Array        DataType = '*'
	BulkString   DataType = '$'
	Error        DataType = '-'
	Integer      DataType = ':'
	SimpleString DataType = '+'
)

type Value struct {
	Type    DataType
	String  string
	Integer int64
	Array   []*Value
	IsNull  bool
}

func readUntilCRLF(rd *bufio.Reader) ([]byte, error) {
	line, err := rd.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read line bytes: %w", err)
	}

	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, errors.New("line not terminated with expected terminator")
	}

	return line[:len(line)-2], nil
}

func deserializeInteger(rd *bufio.Reader) (*Value, error) {
	// Example -> :1234\r\n
	d, err := readUntilCRLF(rd)
	if err != nil {
		return nil, fmt.Errorf("read integer data: %w", err)
	}

	i, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("deserialize integer value: %w", err)
	}

	return &Value{
		Type:    Integer,
		Integer: i,
	}, nil
}

func deserializeSimpleString(rd *bufio.Reader) (*Value, error) {
	// Example -> +Simple string\r\n
	d, err := readUntilCRLF(rd)
	if err != nil {
		return nil, fmt.Errorf("read simple string data: %w", err)
	}

	return &Value{
		Type:   SimpleString,
		String: string(d),
	}, nil
}

func deserializeError(rd *bufio.Reader) (*Value, error) {
	// Example -> -Error message here\r\n
	d, err := readUntilCRLF(rd)
	if err != nil {
		return nil, fmt.Errorf("read error data: %w", err)
	}

	return &Value{
		Type:   Error,
		String: string(d),
	}, nil
}

func deserializeBulkString(rd *bufio.Reader) (*Value, error) {
	// Example -> $5\r\nhello\r\n
	d, err := readUntilCRLF(rd)
	if err != nil {
		return nil, fmt.Errorf("read bulk string data: %w", err)
	}

	strLen, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse string len: %w", err)
	}

	if strLen == -1 {
		return &Value{
			Type:   BulkString,
			IsNull: true,
		}, nil
	}

	strBytes := make([]byte, strLen)
	readLen, err := io.ReadFull(rd, strBytes)
	if err != nil {
		return nil, fmt.Errorf("read bulk string: %w", err)
	}

	if readLen != int(strLen) {
		return nil, fmt.Errorf("short bulk string, expected: %d, got: %d", strLen, readLen)
	}

	crlf := make([]byte, 2)
	n, err := io.ReadFull(rd, crlf)
	if err != nil || n != 2 || crlf[0] != '\r' || crlf[1] != '\n' {
		return nil, fmt.Errorf("bulk string not terminated correctly: %c", crlf)
	}

	return &Value{
		Type:   BulkString,
		String: string(strBytes),
	}, nil
}

func deserializeArray(rd *bufio.Reader) (*Value, error) {
	// Example -> *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	d, err := readUntilCRLF(rd)
	if err != nil {
		return nil, fmt.Errorf("read num elements: %w", err)
	}

	n, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse the num elements: %w", err)
	}

	if n == -1 {
		return &Value{
			Type:   Array,
			IsNull: true,
		}, nil
	}

	arr := make([]*Value, n)

	for i := range n {
		arr[i], err = Deserilize(rd)
		if err != nil {
			return nil, fmt.Errorf("deserialize array element: %w", err)
		}
	}

	return &Value{
		Type:  Array,
		Array: arr,
	}, nil
}

func Deserilize(rd io.Reader) (*Value, error) {
	brd := bufio.NewReader(rd)

	respType, err := brd.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read resp first byte: %w", err)
	}

	switch DataType(respType) {
	case Array:
		return deserializeArray(brd)
	case BulkString:
		return deserializeBulkString(brd)
	case Integer:
		return deserializeInteger(brd)
	case SimpleString:
		return deserializeSimpleString(brd)
	case Error:
		return deserializeError(brd)
	default:
		return nil, fmt.Errorf("unknown resp type: %c", respType)
	}
}
