package main

import (
	"fmt"
	"net"
)

func calcChecksum(body []byte) byte {
	checksum := 0

	for i := 0; i < len(body); i++ {
		checksum += int(body[i])
	}

	return byte(checksum % 256)
}

func createMsg(cmd []byte, payload []byte) []byte {
	hdr := []byte("\xff\xff")
	size := []byte{byte(3 + len(payload))}

	msg := hdr
	msg = append(msg, cmd...)
	msg = append(msg, size...)
	msg = append(msg, payload...)

	cs := calcChecksum(msg[len(hdr):])
	msg = append(msg, cs)
	return msg
}

func sendMsg(msg []byte, ip string, port string) ([]byte, error) {
	conn, _ := net.Dial("tcp", ip+":"+port)
	defer conn.Close()

	_, err := conn.Write(msg)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)

	return buf, err
}

func getWxData(cmd []byte) {
	msg := createMsg(cmd, nil)

	ip := "192.168.4.77"
	port := "45000"

	response, err := sendMsg(msg, ip, port)
	if err != nil {
		fmt.Println(err)
	}

	results, err := parseResponse(response)
	if err != nil {
		fmt.Println(err)
	}
	results.Display()

}

func main() {
	getWxData([]byte("\x27"))
	getWxData([]byte("\x57"))
}
