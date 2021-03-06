=============================================================================     
(同步转异步一般都是采用线程池的方法来实现的异步操作效果,     
go底层采用的是linux内核的线程池做异步操作,很高效的.)     

其实可以这么理解golang中的go关键字:     
  1.go中的函数本身还是正常的同步阻塞函数,     
    所以正常写同步代码时它的效果也就是正常的同步阻塞效果,     
  2.但当你在函数前面加上go关键字后,go会自动的把该函数扔到底层的线程池中去跑,     
    所以自然的有了异步非阻塞的并发运行效果.     
    如果你不需要等待该函数的结果,那么主协程和其他协程都并发运行而已,     
    如果你需要等待该函数的结果,就使用channel在这等待结果即可.     

而且可以这么理解go语言的goroutine+channel并发编程模型:     
  直接想象成你在使用java/python写代码,有一些并发需求,     
  然后你采用'多线程+生产-消费者模式'来写代码实现并发效果而已!     
  只是go的底层使用linux内核的线程池,比一般java/python的多线程更高效而已.     

本质上:     
在其他语言里的[使用多线程+P-C模型来实现异步/并发效果],     
在go里,变成了[使用多协程(底层还是linux内核多线程)+P-C模型来实现异步/并发效果]     
只是所谓的goroutine协程(M:N映射到linux内核的线程)比一般的线程更轻量级更高效.     
=============================================================================

