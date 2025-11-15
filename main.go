package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
)

// PortInfo 1. 定义一个结构体来表示端口状态
type PortInfo struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}
type PortInfoNormal struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}
type PortInfoUsed struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func main() {
	// 初始化日志：同时写到 stdout 和文件 port_scan.log
	logFile, err := os.OpenFile("port_scan.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("无法打开日志文件: %v\n", err)
		return
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("开始端口扫描")

	var usedPort []PortInfoUsed
	var normalPort []PortInfoNormal
	portStatusMap := make(map[int]string)

	for i := 1; i < 49152; i++ {
		func(port int) {
			listenAddr := fmt.Sprintf("localhost:%d", port)

			/**
			udp 实现
			*/

			//udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
			//if err != nil {
			//	println("解析udp地址失败")
			//	return
			//}
			//listener, err := net.ListenUDP("udp", udpAddr)

			/**
			tcp 实现
			*/
			listener, err := net.Listen("tcp", listenAddr)
			if err != nil {
				portStatusMap[port] = err.Error()
				log.Printf("端口 %d 检查失败: %v", port, err)
				return
			}
			portStatusMap[port] = "未占用"
			listener.Close()
		}(i)
	}
	// 2. 将 map 转换为 PortInfo 结构体切片
	for port, status := range portStatusMap {
		if status != "未占用" {
			usedPort = append(usedPort, PortInfoUsed{
				Port:   port,
				Status: status,
			})
		} else {
			normalPort = append(normalPort, PortInfoNormal{
				Port:   port,
				Status: status,
			})
		}
	}

	// 3. 使用 sort.Slice 对切片进行排序
	sort.Slice(usedPort, func(i, j int) bool {
		// 按照 Port 字段从小到大排序
		return usedPort[i].Port < usedPort[j].Port
	})
	sort.Slice(normalPort, func(i, j int) bool {
		// 按照 Port 字段从小到大排序
		return normalPort[i].Port < normalPort[j].Port
	})

	// 4. 将排序后的切片序列化为 JSON
	normalPortJsonData, err := json.MarshalIndent(normalPort, "", "  ")
	if err != nil {
		log.Printf("序列化 normalPort 失败: %v", err)
		return
	}
	usedPortJsonData, err := json.MarshalIndent(usedPort, "", "  ")
	if err != nil {
		log.Printf("序列化 usedPort 失败: %v", err)
		return
	}
	normalPortBuf := []byte(normalPortJsonData)
	usedPortBuf := []byte(usedPortJsonData)

	// 写入文件
	err = os.WriteFile("normal.json", normalPortBuf, 0644)
	if err != nil {
		log.Printf("写入 normal.json 失败: %v", err)
		return
	}
	log.Println("已写入 normal.json", len(normalPort))
	err = os.WriteFile("used_port.json", usedPortBuf, 0644)
	if err != nil {
		log.Printf("写入 used_port.json 失败: %v", err)
		return
	}
	log.Println("已写入 used_port.json", len(usedPort))

	log.Printf("执行完成: 总数量：%d 占用数量: %d 正常数量: %d", len(portStatusMap), len(usedPort), len(normalPort))
}
