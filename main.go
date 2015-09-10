package main

import (
	"flag"
	"fmt"
	"rosa"
	"slingshot/get"
	"slingshot/send"
)

func receiver_mode(dir string, port int) {
	fmt.Println("Slingshot: Waiting file on port", port, "...")
	status := get.Receive(dir, port)
	err := <-status
	if err != nil {
		fmt.Println(err)
	}
}

func sender_mode(file []string, ip string, port int) {
	go send.SendFile(ip, port, file)
	err := <-send.Status
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	portFlag := flag.Int("p", 4242, "A valid port to send or receive file")
	dirFlag := flag.String("o", "./", "output directory")

	flag.Parse()
	tail := flag.Args()

	if len(tail) == 0 {
		receiver_mode(*dirFlag, *portFlag)
	} else if len(tail) > 1 {
		sender_mode(tail[1:], tail[0], *portFlag)
	}
}
