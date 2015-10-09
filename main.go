package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"slingshot/get"
	"slingshot/send"

	"github.com/mrgosti/rosa"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func receiver_mode(dir string, port int, secure bool) {
	if secure {
		fmt.Println("Slingshot secure: Waiting file on port", port, "...")
	} else {
		fmt.Println("Slingshot: Waiting file on port", port, "...")
	}

	status := get.Receive(dir, port, secure)
	err := <-status
	checkErr(err)
}

func sender_mode(ip string, port int, files []string, identity string) {
	go send.SendFile(ip, port, files, identity)
	err := <-send.Status
	checkErr(err)
}

func main() {

	me, err := user.Current()
	checkErr(err)

	rosa.PrivateKeyPath = me.HomeDir + "/.rosa/key.priv"
	rosa.PublicKeyPath = me.HomeDir + "/.rosa/key.pub"
	rosa.FriendListPath = me.HomeDir + "/.rosa/friend_list"

	portFlag := flag.Int("p", 4242, "A valid port to send or receive file")
	dirFlag := flag.String("o", "./", "output directory")
	securFlag := flag.Bool("s", false, "Use secure receive mode (with RoSA)")
	idFlag := flag.String("i", "", "Use public key of this user to send in secure mode (with RoSA)")

	flag.Parse()
	tail := flag.Args()

	if len(tail) == 0 {
		receiver_mode(*dirFlag, *portFlag, *securFlag)
	} else if len(tail) > 1 {
		sender_mode(tail[0], *portFlag, tail[1:], *idFlag)
	}
}
