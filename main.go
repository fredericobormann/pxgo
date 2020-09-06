package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Connecting to server")

	c, err := net.Dial("tcp", "finn-thorben.me:9876")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.Close()

	writer := bufio.NewWriter(c)
	reader := bufio.NewReader(c)

	done := make(chan bool, 1)
	go listen(reader, done)

	for x := 0; x < 800; x++ {
		for y := 0; y < 600; y++ {
			send := fmt.Sprintf("PX %d %d ffffff\n", x, y)
			if _, err := writer.WriteString(send); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err)
	}
	<-done
}

func listen(r *bufio.Reader, done chan<- bool) {
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			done <- true
			break
		}
		fmt.Println(line)
	}
}
