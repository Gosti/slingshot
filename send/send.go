package send

import (
	_ "bufio"
	"errors"
	_ "fmt"
	"io/ioutil"
	"net"
	_ "strconv"
)

var Status chan error

type Pellet struct {
	FileName string
	Content  []byte
}

type Metalic struct {
	P  *Pellet
	CS string
}

func init() {
	Status = make(chan error)
}

func SendFile(IP string, port int, file []string) {
	var ip net.IP

	ip = net.ParseIP(IP)
	if ip == nil {
		ips, err := net.LookupIP(IP)
		if err != nil {
			Status <- errors.New("Slingshot: Not a valid IP")
			return
		}
		ip = ips[0]
	}

	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		Status <- err
		return
	}

	for _, v := range file {
		content, err := ioutil.ReadFile(v)

		if err != nil {
			Status <- err
			return
		}
		conn.Write([]byte(v))
		conn.Write(make([]byte, 1))
		_, err = conn.Write(content)

		if err != nil {
			Status <- err
			return
		}
		conn.Write(make([]byte, 1))
	}

	conn.Close()
	Status <- nil
}
