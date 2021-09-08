// Package controller provides ping method
package controller

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Ping(addr string, pSize, count, timeout, port int) {
	if count == 0 {
		count = math.MaxInt32
	}
	if port > 0 {
		//checkHostPort()
	} else {
		//checkICMP(addr, pSize, count, timeout)
		runPing(addr, pSize, count, timeout)
	}
}

// 定义ICMP数据结构
type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

//
func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

// 检查ICMP
func checkICMP(addr string, pSize, count, timeout int) {
	conn, err := net.Dial("ip4:icmp", addr)
	if err != nil {
		fmt.Println("net dial:", err.Error())
		return
	}
	defer conn.Close()
	
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  1,
		SequenceNum: 1,
	}

	fmt.Printf("PING %s (%s) %d(%d) bytes of data.\n", addr, conn.RemoteAddr(), pSize, pSize+28)

	var buf bytes.Buffer
	// 先在buffer中写入icmp数据报求去校验和
	_ = binary.Write(&buf, binary.BigEndian, icmp)
	//icmp.CheckSum = checkSum(buf.Bytes())
	//// 然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	//buf.Reset()
	//_ = binary.Write(&buf, binary.BigEndian, icmp)
	data := make([]byte, pSize)
	buf.Write(data)
	data = buf.Bytes()

	var totalTime time.Duration
	var SuccessTimes int	// 成功次数
	var FailTimes int		// 失败次数
	var minTime, maxTime time.Duration
	//var mDeviation time.Duration	// 平均方差，放映网速稳定性

	for i:=0; i<count; i++ {
		icmp.SequenceNum = uint16(1)
		// 检验和设为0
		data[2] = byte(0)
		data[3] = byte(0)

		data[6] = byte(icmp.SequenceNum >> 8)
		data[7] = byte(icmp.SequenceNum)
		icmp.CheckSum = checkSum(data)
		data[2] = byte(icmp.CheckSum >> 8)
		data[3] = byte(icmp.CheckSum)

		tStart := time.Now()
		_ = conn.SetDeadline(tStart.Add(time.Duration(timeout) * time.Millisecond))
		n, err := conn.Write(data)
		if err != nil {
			fmt.Println(err.Error())
			FailTimes++
			return
		}
		recv := make([]byte, 65535)
		n, err = conn.Read(recv)
		if err != nil {
			fmt.Println(err.Error())
			FailTimes++
			continue
		}
		costTime := time.Since(tStart)
		if minTime > costTime {
			minTime = costTime
		}
		if maxTime < costTime {
			maxTime = costTime
		}
		totalTime += costTime
		ttl := uint8(recv[8])
		dataLen := len(recv[28:n])
		fmt.Printf("%d bytes from %s: icmp_seq=%d, ttl=%d time=%v\n", dataLen, conn.RemoteAddr(), i+1, ttl, costTime)
		SuccessTimes++
		if (i+1) < count {
			time.Sleep(time.Second)
		}
	}

	fmt.Printf("\n--- %s ping statistics ---\n", addr)
	totalP := SuccessTimes + FailTimes
	loss := float64(FailTimes * 100) / float64(SuccessTimes + FailTimes)
	fmt.Printf("%d packets transmitted, %d received, %.2f%% packet loss, time %v\n",
		totalP, SuccessTimes, loss, totalTime)
	avgTime := float64(int64(totalTime)/int64(totalTime))
	fmt.Printf("rtt min/avg/max = %.3f/%.3f/%.3fms", float64(minTime), avgTime, float64(maxTime))
}

// 检查主机端口
func checkHostPort(addr string) {}

// 调用ping命令
func runPing(addr string, pSize, count, timeout int) {
	cmd := exec.Command("ping",addr, "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeout), "-s", strconv.Itoa(pSize))
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}