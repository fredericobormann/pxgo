package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net"
	"os"
	"sync"
	"time"
)

var url = "finn-thorben.me:9876"
var parallelParts = 80
var width = 800
var height = 600

type sendMode int

const (
	imageMode sendMode = 0
	colorMode sendMode = 1
)

func main() {
	var specifiedMode sendMode
	var img image.Image
	var hexColor string

	if len(os.Args[1:]) >= 2 && os.Args[1] == "--file" {
		specifiedMode = imageMode
		img = openImage(os.Args[2])
	} else if len(os.Args[1:]) >= 2 && os.Args[1] == "--color" {
		specifiedMode = colorMode
		hexColor = os.Args[2]
	} else {
		fmt.Println("No parameter given.")
		return
	}

	fmt.Println("Connecting to server")

	c, err := net.Dial("tcp", url)
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

	if specifiedMode == imageMode {
		for i := 0; i < parallelParts; i++ {
			wg.Add(1)
			fromX := (min(width, img.Bounds().Max.X) / parallelParts) * i
			toX := fromX + (min(width, img.Bounds().Max.X) / parallelParts)
			go sendImagePart(fromX, toX, 0, min(height, img.Bounds().Max.Y), &wg, img)
		}
	} else if specifiedMode == colorMode {
		for i := 0; i < parallelParts; i++ {
			wg.Add(1)
			fromX := (width / parallelParts) * i
			toX := fromX + (width / parallelParts)
			go sendPart(fromX, toX, 0, height, &wg, hexColor)
		}
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

func sendPart(xMin int, xMax int, yMin int, yMax int, wg *sync.WaitGroup, hexColor string) {
	defer wg.Done()

	c, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.Close()

	writer := bufio.NewWriter(c)

	for x := xMin; x < xMax; x++ {
		for y := yMin; y < yMax; y++ {
			send := fmt.Sprintf("PX %d %d %s\n", x, y, hexColor)
			if _, err := writer.WriteString(send); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err)
	}
}

func sendImagePart(xMin int, xMax int, yMin int, yMax int, wg *sync.WaitGroup, img image.Image) {
	defer wg.Done()

	c, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer c.Close()

	writer := bufio.NewWriter(c)

	for x := xMin; x < xMax; x++ {
		for y := yMin; y < yMax; y++ {
			pix := img.At(x, y)
			send := fmt.Sprintf("PX %d %d %s\n", x, y, hexColor(pix))
			if _, err := writer.WriteString(send); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err)
	}
}

func openImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	return img
}

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}
