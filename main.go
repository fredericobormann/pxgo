package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

var parallelParts = 80
var width = 800
var height = 600

func main() {
	fmt.Println("Connecting to server")

	c, err := net.Dial("tcp", "finn-thorben.me:9876")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.Close()

	reader := bufio.NewReader(c)

	done := make(chan bool, 1)
	go listen(reader, done)

	var wg sync.WaitGroup

	for i := 0; i < parallelParts; i++ {
		wg.Add(1)
		go sendPart(i, parallelParts, &wg)
	}

	wg.Wait()
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

func sendPart(part int, numberOfParts int, wg *sync.WaitGroup) {
	defer wg.Done()

	c, err := net.Dial("tcp", "finn-thorben.me:9876")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.Close()

	writer := bufio.NewWriter(c)

	from := (width / numberOfParts) * part
	to := from + (width / numberOfParts)

	for x := from; x < to; x++ {
		for y := 0; y < height; y++ {
			send := fmt.Sprintf("PX %d %d 000ff0\n", x, y)
			if _, err := writer.WriteString(send); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err)
	}
}
