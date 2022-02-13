package gologger

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestPrintln(t *testing.T) {

	thread := 3
	ch := make(chan interface{})
	itempriceLogger := Get("itemprice")
	go func() {
		//1分钟停下
		time.AfterFunc(2*time.Minute, func() {
			ch <- nil
		})
	}()
	for i := 0; i < thread; i++ {
		go func(t int) {
			for {
				msg := fmt.Sprintf("thread %d say hello", t)
				Println(msg)
				itempriceLogger.Println(msg)
				time.Sleep(20 * time.Millisecond)
			}
		}(i)
	}
	<-ch
	Close("", "itemprice")
	log.Println("end")
}

func TestMain(m *testing.M) {
	New(LoggerConf{
		FileDir:    "./log",
		FileName:   "dr",
		Prefix:     "",
		DateFormat: "200601021504",
	})
	New(LoggerConf{
		FileDir:    "./log",
		FileName:   "itemprice",
		Prefix:     "",
		DateFormat: "200601021504",
		Flag:       log.Ldate | log.Ltime,
		Hook: &Hook{
			AfterSplit: func(filepath string) error {
				log.Println("upload file", filepath)
				return nil
			},
		},
	}, "itemprice")
	go StartMonitor()
	os.Exit(m.Run())
}
