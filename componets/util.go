package componets

import (
	"github.com/HelloMrShu/grpc_demo/global"
	"math/rand"
	"net"
	"time"
)

// GenRandStrings 生成随机字符串
func GenRandStrings(max int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"

	bytes := []byte(str)
	var result []byte

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < max; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetIntArrayMax(arr []int) (m int) {
	m = arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}

	return m
}

// GetAvailableAddr 获取可访问ip+port
func GetAvailableAddr() (ip string, port int) {
	ip = GetLocalIp().String()
	port, _ = GetFreePort()

	if ip == "" || port == 0 {
		panic("获取机器ip和端口异常")
	}

	global.ServerConfig.Ip = ip
	global.ServerConfig.Port = port

	return
}

func GetLocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if conn == nil || err != nil {
		global.Logger.Fatal(err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if l == nil || err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
