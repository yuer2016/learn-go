package concurrent

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

/**
* 并发：逻辑上具备同时处理多个任务的能力。
* 并行：物理上在同一时刻执行多个并发任务。
 */

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

/**
* 通道（channel）是显式的，要求操作双方必须知道数据类型和具体通道，并不关心另一端操作身份和数量。
* 可如果另一端未准备妥当，或消息未能及时处理时，会阻塞当前端。
 */
func TestChannel(t *testing.T) {
	data := make(chan int)
	exit := make(chan bool)

	go func() {
		// 循环获取消息，直到通道被关闭
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

/**
* 缓存通道
 */
func TestCacheChannel(t *testing.T) {
	c := make(chan int, 3)

	c <- 10
	c <- 20
	c <- 30

	close(c)

	for i := 0; i < cap(c)+1; i++ {
		x, ok := <-c
		println(i, ":", ok, x)
	}
}

/**
* 通道默认是双向的，并不区分发送和接收端。
* 但某些时候，我们可限制收发操作的方向来获得更严谨的操作逻辑。
* 尽管可用 make 创建单向通道，但那没有任何意义。
* 通常使用类型转换来获取单向通道，并分别赋予操作双方。
 */
func TestSingleChannel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	c := make(chan int)
	var send chan<- int = c
	var recv <-chan int = c

	go func() {
		defer wg.Done()
		// 循环获取消息，直到通道被关闭
		for x := range recv {
			println(x)
		}
	}()

	go func() {
		defer wg.Done()
		defer close(c)

		for i := 0; i < 3; i++ {
			send <- i
		}
	}()

	wg.Wait()
}

/**
* 如要同时处理多个通道，可选用 select 语句。
* 它会随机选择一个可用通道做收发操作
 */
func TestChannelSelect(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	a := make(chan int)
	b := make(chan int)

	// 接收端
	go func() {
		defer wg.Done()

		for {
			var (
				name string
				x    int
				ok   bool
			)

			// 随机选择可用channel接收数据
			select {
			case x, ok = <-a:
				name = "a"
			case x, ok = <-b:
				name = "b"
			}

			// 如果任一通道关闭，则终止接收
			if !ok {
				return
			}

			// 输出接收的数据信息
			println(name, x)
		}
	}()

	// 发送端
	go func() {
		defer wg.Done()
		defer close(a)
		defer close(b)

		for i := 0; i < 10; i++ {
			// 随机选择发送 channel
			select {
			case a <- i:
			case b <- i * 10:
			}
		}
	}()

	wg.Wait()
}

/**
* 要等全部通道消息处理结束 closed，可将已完成通道设置为 nil。
* 这样它就会被阻塞，不再被 select 选中。
 */
func TestChannelSelectNil(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	a := make(chan int)
	b := make(chan int)

	// 接收端
	go func() {
		defer wg.Done()

		for {
			select {
			case x, ok := <-a:
				// 如果通道关闭，则设置为 nil，阻塞
				if !ok {
					a = nil
					break
				}

				println("a", x)
			case x, ok := <-b:
				if !ok {
					b = nil
					break
				}
				println("b", x)
			}

			// 全部结束，退出循环
			if a == nil && b == nil {
				return
			}
		}
	}()

	// 发送端 a
	go func() {
		defer wg.Done()
		defer close(a)

		for i := 0; i < 3; i++ {
			a <- i
		}
	}()

	// 发送端 b
	go func() {
		defer wg.Done()
		defer close(b)

		for i := 0; i < 5; i++ {
			b <- i * 10
		}
	}()

	wg.Wait()
}

/**
* 当所有通道都不可用时，select 会执行 default 语句。
* 这样可避开 select 阻塞，但须注意处理外层循环，以免陷入空耗。
 */
func TestChannelDefaultSelect(t *testing.T) {
	done := make(chan struct{})
	c := make(chan int)

	go func() {
		defer close(done)

		for {
			select {
			case x, ok := <-c:
				if !ok {
					return
				}

				fmt.Println("data:", x)
			default: // 避免 select 阻塞
			}

			fmt.Println(time.Now())
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 5)

	c <- 100
	close(c)

	<-done
}

/**
* 通道本身就是一个并发安全的队列，可用作 ID generator、Pool 等用途。
 */
func TestChannelSemaphore(t *testing.T) {
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup

	// 最多允许2个并发同时执行
	sem := make(chan struct{}, 2)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			//acquire: 获取信号
			sem <- struct{}{}
			//release: 释放信号
			defer func() { <-sem }()

			time.Sleep(time.Second * 2)
			fmt.Println(id, time.Now())
		}(i)
	}

	wg.Wait()
}

/**
* 标准库 time 提供了 timeout 和 tick channel 实现
 */
func TestChannelTime(t *testing.T) {

	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				fmt.Println("timeout...")
				os.Exit(0)
			}
		}
	}()

	go func() {
		tick := time.Tick(time.Second)
		for {
			select {
			case <-tick:
				fmt.Println(time.Now())
			}
		}
	}()

	<-(chan struct{})(nil) // 直接用 nil channel 阻塞进程
}
