// 我们常常需要在后面一个时刻运行 Go 代码，或者在某段时间间隔内重复运行。
// Go 的内置 _定时器_ 和 _打点器_ 特性让这写很容易实现.
// 打点器和定时器的机制有点相似：一个通道用来发送数据
// Go 中最主要的状态管理方式是通过通道间的沟通来完成的但是还是有一些其他的方法来管理状态,
// 例如使用 `sync/atomic`包在多个 Go 协程中进行 _原子计数_ 。
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// ===============================定时器========================
	// 定时器表示在未来某一时刻的独立事件
	timer()
	// ===============================打点器========================
	// _打点器_ 是事件在固定的时间间隔重复执行
	ticker()
	// ===============================工作池========================
	// 使用 Go 协程和通道实现一个_工作池
	// 为了使用 worker 工作池并且收集他们的结果，我们需要2个通道。
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	//3个并行的worker
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	// 发送 9 个 `jobs`
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	// 然后 `close` 这些通道来表示这些就是所有的任务了。
	close(jobs)
	// 最后，收集所有这些任务的返回值。
	for a := 1; a <= 9; a++ {
		<-results
	}
	// ===============================速率限制========================
	// [_速率限制(英)_](http://en.wikipedia.org/wiki/Rate_limiting)
	// 是一个重要的控制服务资源利用和质量的途径
	// Go 通过 Go 协程、通道和[打点器](../tickers/)优美的支持了速率限制。
	rateLimiting()
	// ===============================原子计数器========================
	// 目前我也不理解這個有什麽作用
	atomicCounter()
	// ===============================互斥锁========================
	// 在前面的例子中，我们看到了如何使用原子操作来管理简单的计数器。
	// 对于更加复杂的情况，我们可以使用一个_[互斥锁](http://zh.wikipedia.org/wiki/%E4%BA%92%E6%96%A5%E9%94%81)_
	// 来在 Go 协程间安全的访问数据。
	mutexes()
	// ===============================Go状态协程========================
	// 在前面的例子中，我们用互斥锁进行了明确的锁定来让共享的state 跨多个 Go 协程同步访问。
	// 另一个选择是使用内置的 Go 协程和通道的的同步特性来达到同样的效果。
	// 这个基于通道的方 法和 Go 通过通信以及每个 Go 协程间通过通讯来共享内存，
	// 确保每块数据有单独的 Go 协程所有的思路是一致的。
	statefulGorutines()
}

// =========================定时器========================================
/**
 * @return {[type]} [description]
 */
func timer() {
	// 你告诉定时器需要等待的时间，然后它将提供一个用于通知的通道。
	// 这里的定时器将等待 2 秒。
	timer1 := time.NewTimer(time.Second * 2)
	// `<-timer1.C` 直到这个定时器的通道 `C` 明确的发送了定时器失效的值之前，将一直阻塞。
	<-timer1.C
	fmt.Println("Timer 1 expired")
	// 如果你需要的仅仅是单纯的等待，你需要使用 `time.Sleep`。
	// 定时器是有用原因之一就是你可以在定时器失效之前，取消这个定时器.
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("timer 2 stoped")
	}
}

// =========================打点器========================================
/**
 * 这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
 * @return {[type]} [description]
 */
func ticker() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		// 这里我们在这个通道上使用内置的 `range` 来迭代值每隔500ms 发送一次的值。
		for t := range ticker.C {
			fmt.Println("tick at", t)
		}
	}()
	// 打点器一旦一个打点停止了，将不能再从它的通道中接收到值。
	// 我们将在运行后 1500ms停止这个打点器。
	time.Sleep(time.Millisecond * 1500)
	ticker.Stop()
	fmt.Println("ticker stopped")
}

// =========================工作池========================================
/**
 * 这是我们将要在多个并发实例中支持的任务了。
 * 这些执行者将从 `jobs` 通道接收任务，并且通过 `results` 发送对应的结果。
 * 我们将让每个任务间隔 1s 来模仿一个耗时的任务。
 * @param  {[type]} id   int           [description]
 * @param  {[type]} jobs <-chan        int,          results chan<- int [description]
 * @return {[type]}      [description]
 */
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

// =========================速率限制========================================
/**
 * 限制接收请求的处理，将请求发送给一个相同的通道
 * @return {[type]} [description]
 */
func rateLimiting() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)
	// 这个 `limiter` 通道将每 200ms 接收一个值。这个是速率限制任务中的管理器。
	limiter := time.Tick(time.Millisecond * 200)
	// 通过在每次请求前阻塞 `limiter` 通道的一个接收，我们限制自己每 200ms 执行一次请求。
	for req := range requests {
		// 阻塞200
		<-limiter
		fmt.Println("request", req, time.Now())
	}
	// 有时候我们想临时进行速率限制，并且不影响整体的速率控制
	// 我们可以通过[通道缓冲](channel-buffering.html)来实现。
	// 这个 `burstyLimiter` 通道用来进行 3 次临时的脉冲型速率限制。
	burstyLimiter := make(chan time.Time, 3)
	// 想将通道填充需要临时改变次的值，做好准备。
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}
	// 每 200 ms 我们将添加一个新的值到 `burstyLimiter`中，
	// 直到达到 3 个的限制。
	go func() {
		for t := range time.Tick((time.Millisecond * 200)) {
			burstyLimiter <- t
		}
	}()
	// 现在模拟超过 5 个的接入请求。
	// 它们中刚开始的 3 个将受 `burstyLimiter` 的“脉冲”影响。连续执行
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		// 先将缓存中前三个读取，再每隔200s执行一次
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

// =========================原子计数器========================================
func atomicCounter() {
	// 我们将使用一个无符号整形数来表示（永远是正整数）这个计数器。
	var ops uint64 = 0
	// 为了模拟并发更新，我们启动 50 个 Go 协程，
	// 对计数器每隔 1ms （译者注：应为非准确时间）进行一次加一操作
	for i := 0; i < 50; i++ {
		go func() {
			for {
				// 使用 `AddUint64` 来让计数器自动增加，
				// 使用 `&` 语法来给出 `ops` 的内存地址。
				atomic.AddUint64(&ops, 1)
				// 允许其它 Go 协程的执行
				runtime.Gosched()
			}
		}()
	}
	// 等待一秒，让 ops 的自加操作执行一会。
	time.Sleep(time.Second)
	// 为了在计数器还在被其它 Go 协程更新时，安全的使用它，
	// 我们通过 `LoadUint64` 将当前值得拷贝提取到 `opsFinal` 中。
	// 和上面一样，我们需要给这个函数所取值的内存地址 `&ops`
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops", opsFinal)
}

// =========================互斥锁========================================
/**
 * @return {[type]} [description]
 */
func mutexes() {
	// 在我们的例子中，`state` 是一个 map。
	var state = make(map[int]int)
	// 这里的 `mutex` 将同步对 `state` 的访问。
	var mutex = &sync.Mutex{}
	// `ops` 将记录我们对 state 的操作次数。
	// 爲了比较基于互斥锁的处理方式和我们后面将要看到的其他方式，
	var ops int64 = 0
	// 这里我们运行 100 个 Go 协程来重复读取 state。
	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {
				// 每次循环读取，我们使用一个键来进行访问，
				// `Lock()` 这个 `mutex` 来确保对 `state` 的独占访问，
				// 读取选定的键的值，`Unlock()` 这个mutex，并且 `ops` 值加 1。
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				// 为了确保这个 Go 协程不会再调度中饿死，
				// 我们在每次操作后明确的使用 `runtime.Gosched()`进行释放。
				// 这个释放一般是自动处理的，
				// 像例如每个通道操作后或者 `time.Sleep` 的阻塞调用后相似，
				// 但是在这个例子中我们需要手动的处理。
				runtime.Gosched()
			}
		}()
	}
	// 同样的，我们运行 10 个 Go 协程来模拟写入操作，使用和读取相同的模式。
	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}
	// 让这 10 个 Go 协程对 `state` 和 `mutex` 的操作运行 1 s。
	time.Sleep(time.Second)
	// 获取并输出最终的操作计数。
	opsFinal := atomic.LoadInt64(&ops)
	fmt.Println("ops", opsFinal)
	// 对 `state` 使用一个最终的锁，显示它是如何结束的。
	mutex.Lock()
	fmt.Println("state", state)
	mutex.Unlock()
}

// =========================Go状态协程========================================
type readOp struct {
	key  int
	resp chan int
}
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

/**
 * 在这个例子中，state 将被一个单独的 Go 协程拥有。这就能够保证数据在并行读取时不会混乱。
 * 为了对 state 进行读取或者写入，其他的 Go 协程将发送一条数据到拥有的 Go协程中，
 * 然后接收对应的回复。结构体 `readOp` 和 `writeOp`封装这些请求，并且是拥有 Go 协程响应的一个方式。
 * @return {[type]} [description]
 */
func statefulGorutines() {
	// 和前面一样，我们将计算我们执行操作的次数。
	var ops int64
	// `reads` 和 `writes` 通道分别将被其他 Go 协程用来发布读和写请求。
	reads := make(chan *readOp)
	writes := make(chan *writeOp)
	// 这个就是拥有 `state` 的那个 Go 协程，和前面例子中的map一样，不过这里是被这个状态协程私有的。
	// 这个 Go 协程反复响应到达的请求。
	// 先响应到达的请求，然后返回一个值到响应通道 `resp` 来表示操作成功（或者是 `reads` 中请求的值）
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()
	for w := 0; w < 100; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	opsFinal := atomic.LoadInt64(&ops)
	fmt.Println("ops", opsFinal)
}
