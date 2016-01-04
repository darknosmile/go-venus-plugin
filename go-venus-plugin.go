// go-venus-plug project go-venus-plug.go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-venus-plug/models"
	"net"
	"os"
)

const (
	CURRENT_VERSION = 2
	HEADER_LEN      = 24

	SERIALIZE_JSON = 0x00
	SERIALIZE_BSON = 0x01
	SERIALIZE_JAVA = 0x02

	IS_GZIP = 0x10

	COMMAND_OK        = 0x00000001
	COMMAND_ERROR     = 0xffffffff
	COMMAND_PING      = 0x01000001
	COMMAND_PONG      = 0x01000002
	COMMAND_REQUEST   = 0x02000001
	COMMAND_RESPONSE  = 0x02000002
	COMMAND_HANDSHAKE = 0x03000001
	COMMAND_AUTHEN    = 0x03100000
	COMMAND_NOTIFY    = 0x04000001

	AUTH_DUMMY   = 0x01
	AUTH_USERPWD = 0x02
	AUTH_PKI     = 0x04

	TIMEOUT = 3
)

type Venus struct {
	clientId int32
	conn     *net.TCPConn
}

func Conncect(address string) *Venus {
	addr, err := net.ResolveTCPAddr("tcp", address)
	checkErr(err)
	conn, err := net.DialTCP("tcp", nil, addr)
	checkErr(err)

	data := make([]byte, HEADER_LEN)
	conn.Read(data)

	head := createHead(data)
	fmt.Printf("Head:%+v\n", head)

	//判断是否握手协议
	if fromByte2Int(head.CommondType) != COMMAND_HANDSHAKE {
		fmt.Println("Commond Type ERROR\n")
		os.Exit(1)
	}

	clientId := head.ClientId
	data = make([]byte, head.DateLength)
	conn.Read(data)
	//printData(data)
	shake := createShake(data)

	fmt.Printf("Shake:%+v\n", shake)

	//返回venus连接
	s := &Venus{clientId, conn}
	return s
}

func (s *Venus) Request(apiName, version, jsonData string) (string, string) {
	reserved := make([]byte, 8)
	apiNameByte := []byte(apiName)
	apiNameLength := fromInt32toByte(int32(len(apiNameByte)))
	versionByte := fromInt32toByte(1)
	dataByte := []byte(jsonData)
	dataByteLength := fromInt32toByte(int32(len(dataByte)))
	requestCount := len(reserved) + 4 + len(apiNameByte) + 4 + 4 + len(dataByte) + 16
	//组装报文体
	requestData := make([]byte, requestCount)
	writeData2Byte(&requestData, 0, 8, requestData)
	writeData2Byte(&requestData, 8, 12, apiNameLength)
	lenCur := 12 + len(apiNameByte)
	writeData2Byte(&requestData, 12, lenCur, apiNameByte)
	writeData2Byte(&requestData, lenCur, lenCur+4, versionByte)
	lenCur = lenCur + 4
	writeData2Byte(&requestData, lenCur, lenCur+4, dataByteLength)
	lenCur = lenCur + 4
	writeData2Byte(&requestData, lenCur, lenCur+len(dataByte), dataByte)
	lenCur = lenCur + len(dataByte)
	writeData2Byte(&requestData, lenCur, lenCur+16, make([]byte, 16))

	requestHead := s.makeHead(requestCount, fromInt32toByte(COMMAND_REQUEST))

	msg := make([]byte, HEADER_LEN+requestCount)
	writeData2Byte(&msg, 0, HEADER_LEN, requestHead)
	writeData2Byte(&msg, HEADER_LEN, HEADER_LEN+requestCount, requestData)

	s.conn.Write(msg)

	data := make([]byte, HEADER_LEN)
	s.conn.Read(data)

	//fmt.Printf("DummyPacket Msg:%q\n", data)
	head := createHead(data)
	fmt.Printf("DummyPacket:%+v\n", head)

	if fromByte2Int(head.CommondType) == COMMAND_AUTHEN {
		fmt.Printf("System Error:\n")
		os.Exit(1)
	}

	data = make([]byte, head.PacketLength)
	s.conn.Read(data)

	returnMsg := string(data[4 : len(data)-16])
	traceId := string(data[len(data)-16:])
	fmt.Printf("Return Msg:%+v ,TraceId:%+v", returnMsg, traceId)
	return returnMsg, traceId
}

func (s *Venus) AuthByDummy(userName string) {
	authType := []byte{0x01}
	capabilities := []byte{0x00, 0x00, 0x00, 0x10}
	shakeSerializeType := []byte{0x00}
	client := []byte("VENUS-GO-CLIENT")
	version := []byte("2.0.0-BETA")
	userNameByte := []byte(userName)
	dummyDataCount := len(authType) + len(capabilities) + len(shakeSerializeType) + 4 + len(client) + 4 + len(version) + 4 + len(userNameByte)
	dummyData := make([]byte, dummyDataCount)
	writeData2Byte(&dummyData, 0, 1, authType)
	writeData2Byte(&dummyData, 1, 5, capabilities)
	writeData2Byte(&dummyData, 5, 6, shakeSerializeType)
	writeData2Byte(&dummyData, 6, 10, fromInt32toByte(int32(len(client))))
	lenCur := 10 + len(client)
	writeData2Byte(&dummyData, 10, lenCur, client)
	writeData2Byte(&dummyData, lenCur, lenCur+4, fromInt32toByte(int32(len(version))))
	lenCur = lenCur + 4
	writeData2Byte(&dummyData, lenCur, lenCur+len(version), version)
	lenCur = lenCur + len(version)
	writeData2Byte(&dummyData, lenCur, lenCur+4, fromInt32toByte(int32(len(userNameByte))))
	lenCur = lenCur + 4
	writeData2Byte(&dummyData, lenCur, lenCur+len(userNameByte), userNameByte)

	dummyHead := s.makeHead(dummyDataCount, fromInt32toByte(COMMAND_AUTHEN))
	msg := make([]byte, HEADER_LEN+dummyDataCount)
	writeData2Byte(&msg, 0, HEADER_LEN, dummyHead)
	writeData2Byte(&msg, HEADER_LEN, HEADER_LEN+dummyDataCount, dummyData)
	//fmt.Printf("head length:%v\n", dummyDataCount)
	//fmt.Printf("send data:%q\n", msg)
	s.conn.Write(msg)
	//printData(msg)
	data := make([]byte, 24)
	s.conn.Read(data)
	//printData(data)
	//fmt.Printf("DummyPacket Msg:%q\n", data)
	head := createHead(data)
	fmt.Printf("DummyPacket:%+v\n", head)

	if string(head.CommondType) == string([]byte{0xff, 0xff, 0xff, 0xff}) {
		fmt.Printf("System Error:\n")
		os.Exit(1)
	}

}

func (s *Venus) makeHead(dummyLength int, commondType []byte) []byte {
	head := make([]byte, HEADER_LEN)
	packetLength := fromInt32toByte(int32(dummyLength + HEADER_LEN))
	PacketVersion := []byte{0x00, 0x02}
	serializeType := []byte{0x00}
	flag := []byte{0x00}
	writeData2Byte(&head, 0, 4, packetLength)
	writeData2Byte(&head, 4, 6, PacketVersion)
	writeData2Byte(&head, 6, 10, commondType)
	writeData2Byte(&head, 10, 11, serializeType)
	writeData2Byte(&head, 11, 12, flag)
	writeData2Byte(&head, 12, 16, fromInt32toByte(s.clientId))
	return head
}

func writeData2Byte(dummyData *[]byte, start, end int, addType []byte) {
	j := 0
	for i := start; i < end; i++ {
		(*dummyData)[i] = addType[j]
		j++
	}

}

func createShake(data []byte) (shake *models.PacketHandShake) {
	capabilities := fromByte2Int(data[:4])
	supportAuthenMethod := fromByte2Int(data[4:8])
	randomCharsLen := fromByte2Int(data[8:12])
	challenge := string(data[12 : 12+randomCharsLen])
	//versionCharsLen := fromByte2Int(data[12+randomCharsLen+1 : 16+randomCharsLen])
	version := string(data[16+randomCharsLen:])

	//初始化头对象
	pShake := models.PacketHandShake{capabilities, supportAuthenMethod, challenge, version}

	return &pShake

}

func createHead(data []byte) (head *models.PacketHead) {
	pLength := fromByte2Int(data[:4])
	pVersion := fromByte2Int(data[4:6])
	commondType := data[6:10]
	serializeType := data[10]
	flags := data[11]
	clientId := fromByte2Int(data[12:16])
	clientRequestId := fromByte2Int(data[16:24])
	dateLength := pLength - HEADER_LEN

	//初始化头对象
	pHead := models.PacketHead{pLength, pVersion, commondType, serializeType, flags, clientId, clientRequestId, dateLength}

	return &pHead

}

func fromInt32toByte(num int32) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, num)
	return b_buf.Bytes()
}

func fromByte2Int(data []byte) int32 {
	b_buf := bytes.NewBuffer(data)
	var length int32
	binary.Read(b_buf, binary.BigEndian, &length)
	return length
}

func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1

		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

func printData(data []byte) {
	for _, tmpData := range data {
		fmt.Printf("结果输出:%q\n", ByteToBinaryString(tmpData))
		fmt.Printf("结果输出:%q\n", tmpData&0xff)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
