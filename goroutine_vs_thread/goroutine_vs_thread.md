#### 线程和协程的关系
```
其实线程分2种: 
  内核态线程 + 用户态线程(就是所谓的'协程')

而 协程:线程 有 N:1/1:1/M:N 3种模型,
其中又以go的M:N模型最为复杂和高效,因为它集合了N:1和1:1的优点.
```

```
一个线程什么时候获得CPU时间,能够运行多久,是由OS线程调度器根据某种调度策略所定,这个是系统的线程调度器.
一个协程什么时候获得线程的时间,能够占用线程运行多久,就是Go调度器根据某种调度策略所定,这个是Go调度器.

一个协程的运行,需要占用到线程,一个线程的运行,需要占用到cpu. (G -> M -> P)
可以把协程和线程的关系类比线程和cpu的关系,就是时间片的分配实现并发.(一个线程分配到的cpu时间片再分配给多个协程)

有一点很重要: goroutine是有抢占式调度的.
当一个协程持有一个线程资源后就开始运行,大概在10ms左右后,go调度器就会切换其他协程来使用当前线程资源了(正如系统切换线程一样)
这么做的目的是为了保证所有的协程都能被运行到.
其实就是一个线程所持有的cpu时间片,又被go调度器分配给多个协程了.和系统的线程抢占一样,go的协程也有抢占.
```

```
一直有一个问题:
  那么一个go进程里面到底要开多少个线程在跑着呢?够无数个协程使用吗?那3种需要1:1独占线程的情况呢?
  最终go的高并发效率究竟怎么样?cpu,mem的占用?比如go的http server能处理的http连接数?

```

```
答: 
--------------------------------------------------------------------------------
首先程序里有几种操作需要区分一下:
  1.cpu计算型  ->  通过yield/GoSched这些让出cpu的函数,显然可以做到 N:1
  2.普通的IO阻塞型  ->  例如文件读写,这种IO是阻塞的,需要 1:1
  3.网络IO型  ->  网络IO这种底层的socket这些东西被go特别优化过了,是可以做到 N:1 的！
  4.系统调用型  -> 必须 1:1
  5.调用外部语言函数库的  ->  必须 1:1
在真实的web应用场景中,显然,以 网络IO型 的占比最多!而这也正是go的最适合使用的场景啊!
--------------------------------------------------------------------------------  
在真实引用场景中,那3种需要1:1的情况(普通阻塞IO型/系统调用型/调用外部库)毕竟是少数,
更多的情况则是cpu计算型和网络IO型的操作,所以,go的server的性能和效率应当是不错的.
--------------------------------------------------------------------------------  
```



