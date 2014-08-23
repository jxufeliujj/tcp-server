package socket

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"net"
	"github.com/jxufeliujj/tcp-proto/client"
	"time"
)

var conn *net.TCPConn

const (
	RECV_BUF_LEN = 1024        // set maxium request length to 128KB to prevent flood attack
	MAX_TIME_OUT = time.Minute // set 1 minutes timeout
)

func Start(c *net.TCPConn) {
	conn = c
	conn.SetReadDeadline(time.Now().Add(MAX_TIME_OUT))
	go readLoop()
}

func sendTestData(num int64) {
	msg := new(client.SSC100001)
	msg.Num = proto.Int64(num)
	SendMsg(100001, msg)
}

func readLoop() {
	defer conn.Close() // connection already closed by client
	var request []byte
	for {

		request = make([]byte, RECV_BUF_LEN) // clear last read content
		len, err := conn.Read(request)
		conn.SetReadDeadline(time.Now().Add(MAX_TIME_OUT))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		data := request[:len]

		buf := bytes.NewBuffer(data)
		var l uint32
		var head uint32
		errR := binary.Read(buf, binary.BigEndian, &l)
		checkError(errR)

		errR = binary.Read(buf, binary.BigEndian, &head)
		checkError(errR)

		var b []byte = make([]byte, l)
		errR = binary.Read(buf, binary.BigEndian, &b)
		checkError(errR)
		if head == 100001 {
			msg := new(client.SCS100001)
			msgErr := proto.Unmarshal(b, msg)
			checkError(msgErr)
			fmt.Println("接收解析", head, "->", msg)
		} else {
			fmt.Println("接收", head, "->", string(data))
		}

		sendTestData(int64(time.Now().Second()))
		// errR = binary.Read(buf, binary.BigEndian, &attempts)
		// checkError(errR)
	}
}

func SendMsg(head uint32, msg proto.Message) {
	data, err := proto.Marshal(msg)
	checkError(err)

	var buf bytes.Buffer
	var l uint32 = uint32(len(data))
	err = binary.Write(&buf, binary.BigEndian, l)
	checkError(err)

	err = binary.Write(&buf, binary.BigEndian, head)
	checkError(err)

	err = binary.Write(&buf, binary.BigEndian, data)
	checkError(err)
	fmt.Println("发送", head, "->", buf.String())
	conn.Write(buf.Bytes())
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
