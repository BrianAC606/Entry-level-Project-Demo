package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	timeOut int64
	size    int
	count   int
	icmp    *ICMP = &ICMP{
		Type:        8,
		Code:        0,
		checkSum:    0,
		ID:          0,
		sequenceNum: 0,
	}
)

type ICMP struct {
	Type        uint8
	Code        uint8
	checkSum    uint16
	ID          uint16
	sequenceNum uint16
}

func main() {
	GetCommandArgs()
	//连接绑定ip
	DesIp := os.Args[len(os.Args)-1]
	//raddr, _ := net.ResolveIPAddr("ip", DesIp)
	//laddr := net.IPAddr{IP: net.ParseIP("0.0.0.0")}
	//fmt.Println("raddr:", raddr)
	//fmt.Println("laddr:", laddr)
	//conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	//conn, err := net.DialTimeout("ip4:icmp", raddr.String(), time.Duration(timeOut)*time.Millisecond)
	conn, err := net.DialTimeout("ip4:icmp", DesIp, time.Second*3)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	data := make([]byte, size)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	buffer.Write(data)
	data = buffer.Bytes()

	checkSum := checkSum(data)
	data[2] = byte(checkSum >> 8)
	data[3] = byte(checkSum)

	conn.SetDeadline(time.Now().Add(time.Duration(timeOut) * time.Millisecond))
	_, err = conn.Write(data)
	if err != nil {
		log.Println(err)
	}
	buf := make([]byte, 65535)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(buf[:n])
}

// 设置执行文件的命令行参数
func GetCommandArgs() {
	flag.Int64Var(&timeOut, "t", 1000, "请求超时时长，单位毫秒")
	flag.IntVar(&size, "s", 32, "请求发送的缓冲区大小，单位字节")
	flag.IntVar(&count, "c", 4, "发送请求的次数")
	flag.Parse()
}

// 计算ICMP校验和
func checkSum(data []byte) uint16 {
	len := len(data)
	var sum uint32
	var index int = 0
	for len > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		len -= 2
	}
	if len > 0 {
		sum += uint32(data[index])
	}
	hi16 := sum >> 16
	for hi16 > 0 {
		sum = uint32(int16(sum)) + hi16
		hi16 = sum >> 16
	}
	return 0xFFFF - uint16(sum) //返回sum的补码即校验和
}

//package main
//
//import (
//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"net"
//	"time"
//)
//
//const (
//	MAX_PG = 2000
//)
//
//// 封装 icmp 报头
//type ICMP struct {
//	Type        uint8
//	Code        uint8
//	Checksum    uint16
//	Identifier  uint16
//	SequenceNum uint16
//}
//
//var (
//	originBytes []byte
//)
//
//func init() {
//	originBytes = make([]byte, MAX_PG)
//}
//
//func CheckSum(data []byte) uint16 {
//	len := len(data)
//	var sum uint32
//	var index int = 0
//	for len > 1 {
//		sum += uint32(data[index])<<8 + uint32(data[index+1])
//		index += 2
//		len -= 2
//	}
//	if len > 0 {
//		sum += uint32(data[index])
//	}
//	hi16 := sum >> 16
//	for hi16 > 0 {
//		sum = uint32(int16(sum)) + hi16
//		hi16 = sum >> 16
//	}
//	return 0xFFFF - uint16(sum) //返回sum的补码即校验和
//}
//
//func Ping(domain string, PS, Count int) {
//	var (
//		icmp                      ICMP
//		laddr                     = net.IPAddr{IP: net.ParseIP("0.0.0.0")} // 得到本机的IP地址结构
//		raddr, _                  = net.ResolveIPAddr("ip", domain)        // 解析域名得到 IP 地址结构
//		max_lan, min_lan, avg_lan float64
//	)
//
//	// 返回一个 ip socket
//	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
//
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	defer conn.Close()
//
//	// 初始化 icmp 报文
//	icmp = ICMP{8, 0, 0, 0, 0}
//
//	var buffer bytes.Buffer
//	binary.Write(&buffer, binary.BigEndian, icmp)
//	//fmt.Println(buffer.Bytes())
//	binary.Write(&buffer, binary.BigEndian, originBytes[0:PS])
//	b := buffer.Bytes()
//	binary.BigEndian.PutUint16(b[2:], CheckSum(b))
//
//	//fmt.Println(b)
//	fmt.Printf("\n正在 Ping %s 具有 %d(%d) 字节的数据:\n", raddr.String(), PS, PS+28)
//	recv := make([]byte, 1024)
//	ret_list := []float64{}
//
//	dropPack := 0.0 /*统计丢包的次数，用于计算丢包率*/
//	max_lan = 3000.0
//	min_lan = 0.0
//	avg_lan = 0.0
//
//	for i := Count; i > 0; i-- {
//		/*
//			向目标地址发送二进制报文包
//			如果发送失败就丢包 ++
//		*/
//		if _, err := conn.Write(buffer.Bytes()); err != nil {
//			dropPack++
//			time.Sleep(time.Second)
//			continue
//		}
//		// 否则记录当前得时间
//		t_start := time.Now()
//		conn.SetReadDeadline((time.Now().Add(time.Second * 3)))
//		len, err := conn.Read(recv)
//		/*
//			查目标地址是否返回失败
//			如果返回失败则丢包 ++
//		*/
//		if err != nil {
//			dropPack++
//			time.Sleep(time.Second)
//			continue
//		}
//		t_end := time.Now()
//		dur := float64(t_end.Sub(t_start).Nanoseconds()) / 1e6
//		ret_list = append(ret_list, dur)
//		if dur < max_lan {
//			max_lan = dur
//		}
//		if dur > min_lan {
//			min_lan = dur
//		}
//		fmt.Printf("来自 %s 的回复: 大小 = %d byte 时间 = %.3fms\n", raddr.String(), len, dur)
//		time.Sleep(time.Second)
//	}
//	fmt.Printf("丢包率: %.2f%%\n", dropPack/float64(Count)*100)
//	if len(ret_list) == 0 {
//		avg_lan = 3000.0
//	} else {
//		sum := 0.0
//		for _, n := range ret_list {
//			sum += n
//		}
//		avg_lan = sum / float64(len(ret_list))
//	}
//	fmt.Printf("rtt 最短 = %.3fms 平均 = %.3fms 最长 = %.3fms\n", min_lan, avg_lan, max_lan)
//
//}
//
//func main() {
//	//if len(os.Args) < 3 {
//	//	fmt.Printf("Param domain |data package Sizeof|trace times\n Ex: ./Ping www.so.com 100 4\n")
//	//	os.Exit(1)
//	//}
//	//PS, err := strconv.Atoi(os.Args[2])
//	//if err != nil {
//	//	fmt.Println("you need input correct PackageSizeof(complete int)")
//	//	os.Exit(1)
//	//}
//	//Count, err := strconv.Atoi(os.Args[3])
//	//if err != nil {
//	//	fmt.Println("you need input correct Counts")
//	//	os.Exit(1)
//	//}
//	Ping("www.baidu.com", 48, 5)
//}
