package get

import (
	_ "bufio"
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/mrgosti/rosa"
)

type pellet struct {
	FileName string
	Content  []byte
}

func getPellet(bts []byte) (*pellet, error) {
	var p *pellet

	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func Receive(dir string, port int, secure bool) (status chan error) {
	status = make(chan error)
	fmt.Println(dir, port)
	go getFile(dir, port, secure, status)
	return
}

func save_file(title string, content []byte) {
	title = string(bytes.Trim([]byte(title), "\x00"))
	content = bytes.Trim(content, "\x00")
	file, err := os.Create(title)
	if err != nil {
		fmt.Println("Error during file creation -> ", err)
	}
	_, err = io.WriteString(file, string(content))
	if err != nil {
		fmt.Println("Error during writing -> ", err)
	}
}

func getFile(dir string, port int, secure bool, status chan error) {

	var key *rsa.PrivateKey

	var err error

	if secure {
		key, err = rosa.LoadPrivateKey(rosa.PrivateKeyPath)
		if err != nil {
			status <- err
			return
		}
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		status <- err
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			status <- err
			return
		}
		content := make([]byte, 1)
		go func(c net.Conn) {
			var readed []byte
			defer c.Close()

			for {
				a, _ := c.Read(content)

				if a == 0 && readed != nil {
					fmt.Println(readed)
					if secure {
						readed, err = rosa.Decrypt(readed, key)
						if err != nil {
							status <- err
							return
						}
					}
					p, err := getPellet(readed)
					if err != nil {
						status <- err
						return
					}

					readed = make([]byte, 1)
					save_file(p.FileName, p.Content)
					fmt.Printf("Received => %v\n", p.FileName)
					break
				} else {
					readed = append(readed, content...)
				}
				content = make([]byte, 1)
			}
		}(conn)
	}
	status <- nil
}
