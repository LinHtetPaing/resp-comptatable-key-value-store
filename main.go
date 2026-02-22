package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var aof *Aof

func main() {
	var err error
	aof, err = NewAof("database.aof")
	if err != nil {
		fmt.Println("AOF Error:", err)
		return
	}
	defer aof.Close()

	fmt.Println("Listening on port: 6389")
	listener, err := net.Listen("tcp", ":6389")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	resp := NewResp(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		if value.Typ != "array" || len(value.Array) == 0 {
			continue
		}

		handleCommand(conn, value, true)
	}
}

func handleCommand(conn net.Conn, value Value, shouldWriteToAof bool) {
	command := value.Array[0].Bulk

	fmt.Printf("Received Command: %s\n", command)

	switch command {
	case "PING":
		if conn != nil {
			conn.Write([]byte("+PONG\r\n"))
		}
	case "SET":
		if len(value.Array) != 3 {
			if conn != nil {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
			}
			return
		}
		key := value.Array[1].Bulk
		val := value.Array[2].Bulk

		Set(key, val)

		if shouldWriteToAof {
			aof.Write(value)
		}
		if conn != nil {
			conn.Write([]byte("+OK\r\n"))
		}
	case "GET":
		if len(value.Array) != 2 {
			if conn != nil {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
			}
			return
		}
		key := value.Array[1].Bulk
		val, ok := Get(key)

		if conn != nil {
			if !ok {
				conn.Write([]byte("$-1\r\n")) // NULL
			} else {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			}
		}

	case "DEL":
		if len(value.Array) != 2 {
			if conn != nil {
				conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
			}
			return
		}
		key := value.Array[1].Bulk
		deleted := Del(key)

		if shouldWriteToAof {
			aof.Write(value)
		}
		if conn != nil {
			if deleted {
				conn.Write([]byte(":1\r\n"))
			} else {
				conn.Write([]byte(":0\r\n"))
			}
		}

	case "EXPIRE":
		if len(value.Array) != 3 {
			if conn != nil {
				conn.Write([]byte("-ERR wrong number of arguments for 'expire' command\r\n"))
			}
			return
		}
		key := value.Array[1].Bulk
		seconds, _ := strconv.Atoi(value.Array[2].Bulk)

		success := Expire(key, time.Duration(seconds)*time.Second)

		if shouldWriteToAof {
			aof.Write(value)
		}
		if conn != nil {
			if success {
				conn.Write([]byte(":1\r\n"))
			} else {
				conn.Write([]byte(":0\r\n"))
			}
		}

	default:
		if conn != nil {
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}

}
