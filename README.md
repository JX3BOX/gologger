golang 日志切割

## 快速开始

```go
package main
import "jx3box.com/JX3Box/gologger"

func main(){
    gologger.New(gologger.LoggerConf{
		FileDir:    "./log",
		FileName:   "dr",
		Prefix:     "test",
		DateFormat: "2006010215",
	})
    gologger.Println("1111")
}
```

### LoggerConf

```
FileDir: 日志文件存储地址
FileName： 日志文件
Prefix    每行文本前缀文本
Flag      日志flag，如：log.Ldate | log.Ltime
DateFormat: 切割文件时间粒度。 默认为20060102，20060102表示切割到天，2006010215  表示切割到小时
Hook  钩子处理文件切割事件
```


### Hook

```golang
type Hook struct {
	// 文件切割完成后到调用此函数
	// @params filepath 文件路径
	AfterSplit func(filepath string) error
}
```