package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

type PortStruct struct {
	Port int
	Open bool
}

var reqChan chan int
var ansChan chan PortStruct

func main() {
	color.White("Enter the link")
	color.New(color.FgHiBlack).Println("example: \"something.com\"")

	scanner := bufio.NewScanner(os.Stdin)
	var input string
	if scanner.Scan() {
		input = scanner.Text()
	}

	reqChan = make(chan int, 65535)
	ansChan = make(chan PortStruct)

	var wg sync.WaitGroup

	for i := 1; i <= 65535; i++ {
		reqChan <- i
	}
	close(reqChan)

	workerCount := 1000
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go CheckPort(input, reqChan, ansChan, &wg)
	}

	go func() {
		for i := 1; i <= 65535; i++ {
			val := <-ansChan
			if !val.Open {
				color.Red("port %d is closed", val.Port)
			} else {
				color.Green("port %d is open", val.Port)
			}
		}
	}()

	wg.Wait()
}

func CheckPort(input string, reqChan chan int, ansCh chan PortStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	for reqPort := range reqChan {
		port := fmt.Sprintf("%s:%d", input, reqPort)
		ans, err := net.DialTimeout("tcp", port, 10*time.Second)
		if err != nil {
			ansCh <- PortStruct{
				Port: reqPort,
				Open: false,
			}
		} else {
			ansCh <- PortStruct{
				Port: reqPort,
				Open: true,
			}
			ans.Close()
		}
	}
}
