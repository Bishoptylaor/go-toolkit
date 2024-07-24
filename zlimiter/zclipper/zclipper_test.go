package zclipper

import (
	"fmt"
	"testing"
	"time"
)

func TestIntervalFunc(t *testing.T) {
	var ch chan bool
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("end")
				return
			case <-time.NewTicker(time.Second).C:
				fmt.Println("interval calling at ", time.Now().Unix())
			}
		}
	}()

	time.Sleep(time.Second * 15)
	ch <- true
}
