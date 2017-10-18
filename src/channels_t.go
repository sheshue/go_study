/**
 * _通道_ 是连接多个 Go 协程的管道。你可以从一个 Go 协程
 * 将值发送到通道，然后在别的 Go 协程中接收。
 */
package main

import (
	"fmt"
	"time"
)

func main() {
	// ===============================通道========================
	// 使用 `make(chan val-type)` 创建一个新的通道。
	// 通道类型就是他们需要传递值的类型。
	messages := make(chan string)
	go func() {
		// 使用 `channel <-` 语法 _发送_ 一个新的值到通道中。这里
		// 我们在一个新的 Go 协程中发送 `"ping"` 到上面创建的
		// `messages` 通道中。
		messages <- "ping"
	}()
	// 使用 `<-channel` 语法从通道中 _接收_ 一个值.
	msg := <-messages
	// 将接收我们在上面发送的 `"ping"` 消息并打印出来。
	fmt.Println(msg)
	//========================通道方向===============================
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
	// ==========================通道缓冲=============================
	// 默认通道是 _无缓冲_ 的，这意味着只有在对应的接收通道（`<- chan`）
	// 准备好接收时，才允许进行发送值到通道（`chan <-`）。
	// _可缓存通道_ 允许在没有对应接收方的情况下，缓存限定数量的值
	// 这里我们 `make` 了一个通道，最多允许缓存 2 个值。
	chan2 := make(chan string, 2)
	// 因为这个通道是有缓冲区的，即使没有一个对应的并发接收方，我们仍然可以发送这些值。
	chan2 <- "buffered"
	chan2 <- "channel"
	// 然后我们可以像前面一样接收这两个值。
	fmt.Println(<-chan2)
	fmt.Println(<-chan2)
	// ==========================通道同步=============================
	// 我们可以使用通道来同步 Go 协程间的执行状态。
	// 这里是一个使用阻塞的接受方式来等待一个 Go 协程的运行结束。
	done := make(chan bool, 1)
	go worker(done)
	// 程序将在接收到通道中 worker 发出的通知前一直阻塞。
	<-done
	// ==========================通道选择器=============================
	// Go 的_通道选择器_ 让你可以同时等待多个通道操作。
	// Go 协程和通道以及选择器的结合是 Go 的一个强大特性。
	chanSelector()
}

// =========================通道方向========================================
/**
 * 当使用通道作为函数的参数时，你可以指定这个通道是不是只用来发送或者接收值。
 * 这个特性提升了程序的类型安全性。
 * `ping` 函数定义了一个只允许发送数据的通道。尝试使用这个通道来接收数据将会得到一个编译时错误。
 */
func ping(pings chan<- string, msg string) {
	pings <- msg
}

/**
 * `pong` 函数允许通道（`pings`）来接收数据，另一通道，（`pongs`）来发送数据。
 * @param  {[type]} pings <-chan        string, pongs chan<- string [description]
 * @return {[type]}       [description]
 */
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

// =========================通道同步========================================
/**
 * 这是一个我们将要在 Go 协程中运行的函数。
 * `done` 通道将被用于通知其他 Go 协程这个函数已经工作完毕。
 * @param  {[type]} done chan          bool [description]
 * @return {[type]}      [description]
 */
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

// =========================通道选择器========================================
/**
 * 在我们的例子中，我们将从两个通道中选择。
 * @return {[type]} [description]
 */
func chanSelector() {
	c1 := make(chan string)
	c2 := make(chan string)
	// 各个通道将在若干时间后接收一个值，这个用来模拟例如并行的 Go 协程中阻塞的 RPC 操作
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "two"
	}()
	// 我们使用 `select` 关键字来同时等待这两个值，并打印各自接收到的值。
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("reveived", msg2)
		}
	}
}
