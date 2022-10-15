golang 日志

## 快速开始

### 控制台日志

```go
package main
import "gopkg.in/JX3Box/gologger.v2"

func main(){
	gologger.InitLogger(gologger.DebugLevel)
    gologger.Println("1111")
}
```

### 写入文件

```go
package main
import (
	"gopkg.in/JX3Box/gologger.v2/rotatefile"
	"gopkg.in/huyinghuan/lumberjack.v4"
)

func main(){
	l, err := lumberjack.NewRoller(
		"/var/log/myapp/foo.log",
		&lumberjack.Options{
			Filename:   "/var/log/myapp/foo.log",
			MaxAge:     28, //days
			RotateType: RotateDaily, //optioanl, RotateHourly or RotateDaily, If not set, use rotate by size
			RotateTime: 1, // optional, default 1
		})
	if err!=nil{
		panic(err)
	}
	rotatefile.New(l)
    gologger.Println("1111")
}

```