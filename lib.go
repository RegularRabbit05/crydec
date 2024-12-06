package main

import "C"
import (
	"encoding/binary"
	"encoding/json"
	"net"
	"unsafe"
)

func main() {
}

//export sttJsonToText
func sttJsonToText(decBuf *uint8, decBufLen uint32, js string) {
	strRes := unsafe.Slice(decBuf, decBufLen)
	strRes[0] = 0
	if len(js) == 0 {
		return
	}
	type Tp struct {
		Text string `json:"text"`
	}
	var tp Tp
	err := json.Unmarshal([]byte(js), &tp)
	p := tp.Text
	if err != nil {
		type Tp2 struct {
			Text string `json:"partial"`
		}
		var tp2 Tp2
		err = json.Unmarshal([]byte(js), &tp2)
		if err != nil {
			return
		}
		p = tp2.Text
	}
	var i uint32
	for i = 0; i < uint32(len(p)) && i < decBufLen-1; i++ {
		strRes[i] = p[i]
	}
	strRes[i] = 0
}

var connection net.Conn

//export connectToApp
func connectToApp() {
	conn, err := net.Dial("tcp", "localhost:56864")
	if err != nil {
		panic(err)
	}
	connection = conn
}

//export closeConnection
func closeConnection() {
	connection.Close()
}

//export sendStringToApp
func sendStringToApp(str string) {
	var num uint16 = uint16(len(str))
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, num)
	buf = append(buf, []byte(str)...)
	connection.Write(buf)
}
