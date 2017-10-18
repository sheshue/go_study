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
	// ===============================通道实例========================
	channels()
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
	// ==========================超时处理=============================
	// _超时_
	// 对于一个连接外部资源，或者其它一些需要花费执行时间的操作的程序而言是很重要的。
	// 得益于通道和 `select`，在 Go中实现超时操作是简洁而优雅的。
	timeouts()
	// ==========================非阻塞通道操作=============================
	// 常规的通过通道发送和接收数据是阻塞的。
	// 然而，我们可以使用带一个 `default` 子句的 `select`
	// 来实现_非阻塞_ 的发送、接收，
	// 甚至是非阻塞的多路 `select`(意思是select 嵌套？)。
	nonBlockingChannel()
	// ==========================通道的关闭=============================
	// _关闭_ 一个通道意味着不能再向这个通道发送值了。
	// 这个特性可以用来给这个通道的接收方传达工作已将完成的信息。
	done2 := make(chan bool, 1)
	closeChannel(done2)
	<-done2
}

// =========================通道实例========================================
/**
 * 通道实例
 * @return {[type]} [description]
 */
func channels() {
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

// =========================超时处理========================================、
/**
 * 使用`select` 超时方式，需要使用通道传递结果。
 * 这对于一般情况是个好的方式，因为其他重要的 Go 特性是基于通道和`select` 的。
 * @return {[type]} [description]
 */
func timeouts() {
	// 假如我们执行一个外部调用，并在 2 秒后通过通道 `c1` 返回它的执行结果。
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "result 1"
	}()
	// 这里是使用 `select` 实现一个超时操作。
	// `res := <- c1` 等待结果，`<-Time.After` 等待超时 时间 1 秒后发送的值。
	// 由于 `select` 默认处理第一个已准备好的接收操作，
	// 如果这个操作超过了允许的 1 秒的话，将会执行超时case。
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}

	// 如果我允许一个长一点的超时时间 3 秒，将会成功的从 `c2`
	// 接收到值，并且打印出结果。
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout 2")
	}
}

// =========================非阻塞通道操作========================================
/**
 * 非阻塞通道接受，发送实例
 * @return {[type]} [description]
 */
func nonBlockingChannel() {
	messages := make(chan string)
	signals := make(chan bool)
	// 这里是一个非阻塞`接收`的例子。
	// 如果在 `messages` 中存在，然后 `select` 将这个值带入 `<-messages` `case`。
	// 如果不是，就直接到 `default` 分支中。
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}
	// 一个非阻塞`发送`的实现方法。
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}
	// 我们可以在 `default` 前使用多个 `case` 子句来实现一个多路的非阻塞的选择器。
	// 这里我们试图在 `messages`和 `signals` 上同时使用非阻塞的接受操作。
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

// =========================通道的关闭========================================
func closeChannel(done chan bool) {
	jobs := make(chan int, 5)
	// 这是工作 Go 协程。
	// 使用 `j, more := <- jobs` 循环的从 `jobs` 接收数据。
	// 在接收的这个特殊的二值形式的值中，
	// 如果 `jobs` 已经关闭了，并且通道中所有的值都已经接收完毕, 那么 `more` 的值将是 `false`。
	// 当我们完成所有的任务时，将使用这个特性通过 `done` 通道去进行通知。
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				//发送一个值通知已经完成
				done <- true
				return
			}
		}
	}()
	// 这里使用 `jobs` 发送 3 个任务到工作函数中，然后
	// 关闭 `jobs`。
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")
	// 我们使用前面学到的[通道同步]
	// 方法等待任务结束。
}
