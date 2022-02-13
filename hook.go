package gologger

type Hook struct {
	// 文件切割完成后到调用此函数
	AfterSplit func(filepath string) error
}
