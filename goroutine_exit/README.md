## 多个goroutine的退出控制

##### 我们通过select和default分支可以很容易实现一个Goroutine的退出控制:
```go
func worker(cannel chan bool) {
	for {
		select {
		default:
			fmt.Println("hello")
			// 正常工作
		case <-cannel:
			// 退出
		}
	}
}

func main() {
	cannel := make(chan bool)  // 无缓冲channel
	go worker(cannel)

	time.Sleep(time.Second)
	cannel <- true    // 发送一个退出信号,则一个worker会退出
}
```

##### 但是管道的发送操作和接收操作是一一对应的，如果要停止多个Goroutine那么可能需要创建同样数量的管道，这个代价太大了。
##### 其实我们可以通过close关闭一个管道来实现 [广播] 的效果，所有从关闭管道接收的操作均会收到一个零值和一个可选的失败标志。

```go
func worker(cannel chan bool) {
	for {
		select {
		default:
			fmt.Println("hello")
			// 正常工作
		case <-cannel:
			// 退出
		}
	}
}

func main() {
	cancel := make(chan bool)

	for i := 0; i < 10; i++ {
		go worker(cancel)
	}

	time.Sleep(time.Second)
	close(cancel)  // 关闭了channel后,所有的worker里面都会获取到零值,相当于给所有的worker发送了退出信号
}
```

##### 目前main线程并没有等待各个工作Goroutine退出工作完成, 我们可以结合sync.WaitGroup来改进:
```go
func worker(wg *sync.WaitGroup, cannel chan bool) {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-cannel:
			return
		}
	}
}

func main() {
	cancel := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(&wg, cancel)
	}

	time.Sleep(time.Second)
	close(cancel)
	wg.Wait()
}
```


##### 在Go1.7发布时，标准库增加了一个context包，用来简化对于处理单个请求的多个Goroutine之间与请求域的数据、超时和退出等操作
##### 当并发体超时或main主动停止工作者Goroutine时，每个工作者都可以安全退出。
##### 我们可以用context包来重新实现前面的线程安全退出或超时的控制: 

```go
func worker(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}

	time.Sleep(time.Second)
	cancel()

	wg.Wait()
}
```




