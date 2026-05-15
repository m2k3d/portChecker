package main

import (
	"fmt"
	"sync"
	"testing"
)

func Test_CheckPort(t *testing.T) {
	reqChan = make(chan int, 1)
	ansChan = make(chan PortStruct, 1)
	var wg sync.WaitGroup

	reqChan <- 80
	close(reqChan)

	wg.Add(1)
	go CheckPort("localhost", reqChan, ansChan, &wg)

	wg.Wait()

	testResult := <-ansChan
	if testResult.Open {
		fmt.Println("test is completed, port is open")
	} else {
		fmt.Println("test is completed, port is closed")
	}
}
