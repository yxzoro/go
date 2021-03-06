```
我们还可以进一步地通过类型查询或类型断言来获取底层真实的错误类型，这样就可以获取更详细的错误信息。
不过一般情况下我们并不关心错误在底层的表达方式，我们只需要知道它是一个错误就可以了。
当返回的错误值不是nil时，我们可以通过调用error接口类型的Error方法来获得字符串类型的错误信息。

在Go语言中，错误被认为是一种可以预期的结果；而异常则是一种非预期的结果，发生异常可能表示程序中存在BUG或发生了其它不可控的问题。
Go语言推荐使用recover函数将内部异常转为错误处理，这使得用户可以真正的关心业务相关的错误处理。

如果某个接口简单地将所有普通的错误当做异常抛出，将会使错误信息杂乱且没有价值。
就像在main函数中直接捕获全部一样，是没有意义的：

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	...
}
捕获异常不是最终的目的。如果异常不可预测，直接输出异常信息是最好的处理方式。


前文我们说到，Go语言中的导出函数一般不抛出异常，一个未受控的异常可以看作是程序的BUG。
但是对于那些提供类似Web服务的框架而言；它们经常需要接入第三方的中间件。因为第三方的中间件是否存在BUG是否会抛出异常，Web框架本身是不能确定的。
为了提高系统的稳定性，Web框架一般会通过recover来防御性地捕获所有处理流程中可能产生的异常，然后将异常转为普通的错误返回。

让我们以JSON解析器为例，说明recover的使用场景。考虑到JSON解析器的复杂性，即使某个语言解析器目前工作正常，也无法肯定它没有漏洞。
因此，当某个异常出现时，我们不会选择让解析器崩溃，而是会将panic异常当作普通的解析错误，并附加额外信息提醒用户报告此错误。

func ParseJSON(input string) (s *Syntax, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("JSON: internal error: %v", p)       // 最外层的异常捕捉
		}
	}()
	// ...parser...
}
标准库中的json包，在内部递归解析JSON数据的时候如果遇到错误，
会通过抛出异常的方式来快速跳出深度嵌套的函数调用，然后由最外一级的接口通过recover捕获panic，然后返回相应的错误信息。

Go语言库的实现习惯: 即使在包内部使用了panic，但是在导出函数时会被转化为明确的错误值。
```

[go的错误处理策略](https://github.com/chai2010/advanced-go-programming-book/blob/master/ch1-basic/ch1-07-error-and-panic.md)



