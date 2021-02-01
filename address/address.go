package address

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ()  {
	
}


var usedPort = map[int]struct{}{}
var (
	min = 10000
	max = 11000
)

func ResetAddress(address string) {
	s := strings.Split(address, ":")
	if len(s) == 2 {
		port, err := strconv.Atoi(s[1])
		if err == nil {
			ResetPort(port)
		}
	}
}

func ResetPort(port int) {
	fmt.Println("reset", port)
	delete(usedPort, port)
}

func GetFreePort() (int, error) {
	for port := min; port <= max; port++ {
		if _, ok := usedPort[port]; ok {
			continue
		}

		address := fmt.Sprintf("localhost:%d", port)
		_, err := getPort(address)
		fmt.Println("set", port)
		usedPort[port] = struct{}{}

		if err == nil {
			return port, nil
		}
	}
	return 0, fmt.Errorf("free port not found")

}

func getPort(address string) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}


func d(){
	net.
}