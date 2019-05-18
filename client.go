package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		status, _ := bufio.NewReader(conn).ReadString('\r')
		fmt.Print(status)
		if status == "bye!\n\r" {
			break
		}
		if scanner.Scan() {
			str := scanner.Text()
			fmt.Fprintf(conn, str + "\n")
		}
	}
}
