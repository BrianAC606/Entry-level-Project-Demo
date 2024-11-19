package main

import (
	"ProjectDemo/DNSDemo/dns_body"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	Domain := os.Args[len(os.Args)-1]
	DNSLocalServer := "202.202.32.33:53"

	reqHeader := dns_body.NewHeader()
	reqHeader.SetFlags(0, 0, 0, 0, 0, 0, 0)

	dnsRequest := dns_body.NewQueries(Domain)
	//fmt.Println(dnsRequest.GetQueriesBtyes())

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, reqHeader)
	binary.Write(&buffer, binary.BigEndian, dnsRequest.GetQueriesBtyes())

	conn, err := net.Dial("udp", DNSLocalServer)
	if err != nil {
		log.Fatalln("connect error:", err)
	}
	defer conn.Close()
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Fatalln("write error:", err)
	}

	dnsResponse := make([]byte, 1024)
	_, err = conn.Read(dnsResponse)
	if err != nil {
		log.Fatalln("read error:", err)
	}

	var (
		resHeader = dnsResponse[:12]
		resData   = dnsResponse[12:]
	)
	index := 0

	var (
		queryNum  = uint16(resHeader[4]<<8) + uint16(resHeader[5])
		answerNum = uint16(resHeader[6]<<8) + uint16(resHeader[7])
	)

	resQueries := make([][]byte, queryNum)
	for i := 0; i < int(queryNum); i++ {
		for ; resData[index] != 0; index += int(resData[index]) + 1 {
			resQueries[i] = append(resQueries[i], resData[index:index+int(resData[index])+1]...)
		}
		resQueries[i] = append(resQueries[i], resData[index:index+5]...)
		index += 5
	}

	var desIp [][]byte
	for i := 0; i < int(answerNum); i++ {
		Tye := resData[index+3]
		index += 2 + 2 + 2 + 4
		dataLen := uint16(resData[index])<<8 + uint16(resData[index+1])
		index += 2
		if Tye == 1 {
			ip := resData[index : index+4]
			desIp = append(desIp, ip)
		}
		index += int(dataLen)
	}

	fmt.Printf("正在解析域名:%s\n", Domain)
	for _, ip := range desIp {
		fmt.Printf("该域名的ip为: %d.%d.%d.%d\n", ip[0], ip[1], ip[2], ip[3])
	}
}
