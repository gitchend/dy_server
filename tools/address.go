package tools

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func AddressSplit(addr string) (ip, port string) {
	strs := strings.Split(addr, ":")
	return strs[0], strs[1]
}

func AddressMerge(ip, port string) (addr string) {
	return fmt.Sprintf("%s:%s", ip, port)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
