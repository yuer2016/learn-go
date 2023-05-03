package concurrent

import (
	"math"
	"runtime"
	"sync"
	"testing"
	"time"
)

func sum(id uint16) {
	var sum uint32

	for i := 0; i < math.MaxInt16; i++ {
		sum += uint32(i)
	}

	println(sum, id)
}

// test wait group
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	for i := 0; i < 2; i++ {

		go func(id int) {
			defer wg.Done()
			sum(uint16(id))
		}(i)

	}

	wg.Wait()
}

/*
* runtime.Goexit() can be terminated goroutine
 */
func TestExitGoroutine(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer println("defer.A")

		func() {
			defer println("defer.B")
			runtime.Goexit()
			println("B")
		}()

		println("A")
	}()

	wg.Wait()
}

func TestGosched(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < 6; i++ {
			println(i)
			if i == 3 {
				runtime.Gosched()
			}
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(10) * time.Second)
		println("Gosched like yeild operation")
	}()

	wg.Wait()
}

func TestChannel(t *testing.T) {
	data := make(chan int)
	exit := make(chan bool)

	go func() {
		for d := range data {
			println(d)
		}

		println("receive over.")
		exit <- true
	}()

	data <- 1
	data <- 2
	data <- 3
	close(data)
	println("send over.")
	<-exit
}
