package send

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/mrgosti/rosa"
)

var Status chan error

type pellet struct {
	FileName string
	Content  []byte
}

func (p *pellet) toBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendFile(IP string, port int, files []string, identity string) {
	rosa.LoadFriends(rosa.FriendListPath)
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

	for _, filename := range files {
		content, err := ioutil.ReadFile(filename)

		if err != nil {
			Status <- err
			return
		}

		p := &pellet{filename, content}

		tosend, err := p.toBytes()
		fmt.Println(tosend)
		if err != nil {
			Status <- err
			return
		}

		if identity != "" {
			f := rosa.SeekByName(identity)
			if f == nil {
				Status <- errors.New("Slingshot: Not a valid identity")
				return
			}
			tosend, err = f.Encrypt(tosend)
			fmt.Println(tosend)
			if err != nil {
				Status <- err
				return
			}
		}

		_, err = conn.Write(tosend)

		if err != nil {
			Status <- err
			return
		}
	}

	conn.Close()
	Status <- nil
}

func init() {
	Status = make(chan error)
}
