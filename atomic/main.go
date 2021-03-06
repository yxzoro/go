package main

import (
	"time"
	"log"
	"sync/atomic"
)


func main() {
	log.Println("test atomic, i = 0")
	var i int64 = 0
	for x := 0; x < 100; x++ {
		go write(&i, x)
		go read(&i, x)
	}

	time.Sleep(2*time.Second)
	log.Println("end, i = ", i)
}

func write(i *int64, x int) {
	log.Println(x, " add 1 to i")
	atomic.AddInt64(i, 1)  // atomic write 
	log.Println(x, " add 1 to i ok")
}

func read(i *int64, x int) {
	log.Println(x, " read i = ", atomic.LoadInt64(i))  // atomic read
}

// atomic make sure that: i = 100
// it's just like add a lock when atomic modify i




