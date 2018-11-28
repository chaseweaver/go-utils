package main

import (
	"fmt"
	"net"
	"strings"
)

func init() {
	cmd := Command{
		Name:        "ip",
		Description: "Lists public IP address",
		Action:      publicIP,
	}
	Commands[cmd.Name] = cmd
}

func publicIP() {
	fmt.Println("Public IP: ", outboundIP())
}

// outboundIP returns public IP address
func outboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
