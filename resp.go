package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

// reading bytes and then convert into value
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()

	case BULK:
		return r.readBulk()

	default:
		fmt.Printf("Unknown type: %v\n", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{Typ: "array"}

	lenLine, _, err := r.reader.ReadLine()
	if err != nil {
		return v, err
	}

	length, err := strconv.Atoi(string(lenLine))
	if err != nil {
		return v, err
	}

	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{Typ: "bulk"}

	lenLine, _, err := r.reader.ReadLine()
	if err != nil {
		return v, err
	}

	length, err := strconv.Atoi(string(lenLine))
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	r.readLine()

	return v, nil
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	default:
		return []byte("")
	}
}

func (v Value) marshalArray() []byte {
	res := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
	for _, val := range v.Array {
		res = append(res, val.Marshal()...)
	}
	return res
}

func (v Value) marshalBulk() []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Bulk), v.Bulk))
}
