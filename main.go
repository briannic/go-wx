package main

import (
	"fmt"
	"net"
)

func sum(arr []byte) int {
	result := 0
	for _, v := range arr {
		result += int(v)
	}
	return result
}

func calcChecksum(body [][]byte) byte {
	checksum := 0

	for i := 0; i < len(body); i++ {
		checksum += sum(body[i])
	}

	return byte(checksum % 256)
}

func createMsg() []byte {
	hdr := []byte("\xff\xff")
	cmd := []byte("\x27")
	pay := []byte(nil)
	size := []byte{byte(3 + len(pay))}

	cs := calcChecksum([][]byte{cmd, pay, size})

	msg := hdr
	msg = append(msg, cmd...)
	msg = append(msg, size...)
	msg = append(msg, pay...)
	msg = append(msg, cs)
	return msg
}

func sendMsg(msg []byte) ([]byte, error) {
	ip := "192.168.4.77:45000"
	conn, _ := net.Dial("tcp", ip)
	defer conn.Close()

	_, err := conn.Write(msg)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)

	return buf, err
}

func main() {
	msg := createMsg()
	response, _ := sendMsg(msg)
	err := parseResponse(response)
	if err != nil {
		fmt.Println(err)
	}
}
