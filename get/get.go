package get

import (
	_ "bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func Receive(dir string, port int) (status chan error) {
	status = make(chan error)
	fmt.Println(dir, port)
	go getFile(dir, port, status)
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

func getFile(dir string, port int, status chan error) {
	var readed []byte
	var title string = ""
	var fileContent []byte

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
			defer c.Close()

			for {
				a, b := c.Read(content)
				if a == 0 || b != nil {
					break
				}
				readed = append(readed, content...)
				if content[0] == 0 {
					if title == "" {
						title = string(readed)
						readed = readed[:cap(readed)]
					} else {
						fileContent = readed
						readed = readed[:cap(readed)]
						fmt.Printf("Received => %v\n", title)
						save_file(title, fileContent)
						title = ""
					}
				}
				content = make([]byte, 1)
			}
		}(conn)
	}
	status <- nil
}
