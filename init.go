package gologger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// const DATE_FORMAT = "20060102"
//var fileLog *FileLogger

var multLog = make(map[string]*FileLogger)
var multLogMut = new(sync.RWMutex)

type LoggerConf struct {
	FileDir    string `yaml:"fileDir"`
	FileName   string `yaml:"fileName"`
	Prefix     string `yaml:"prefix"`
	Flag       int    `yaml:"flag"`
	DateFormat string `yaml:"dateFormat"`
	Hook       *Hook  `yaml:"-"`
}

type FileLogger struct {
	fileDir     string
	fileName    string
	prefix      string
	logFlag     int
	date        *time.Time
	logFile     *os.File
	logFilepath string
	lg          *log.Logger
	mu          *sync.RWMutex
	logChan     chan string
	dateFormat  string
	hook        *Hook
}

func (f *FileLogger) getLogFilePath() string {
	f.logFilepath = filepath.Join(f.fileDir, f.fileName) + "." + f.date.Format(f.dateFormat) + ".log"
	return f.logFilepath
}

// 日志文件是否必须分割
func (f *FileLogger) isMustSplit() bool {
	t, _ := time.Parse(f.dateFormat, time.Now().Format(f.dateFormat))
	return t.After(*f.date)
}

// 日志文件是否存在，不存在则创建
func (f *FileLogger) isExistOrCreate() {
	_, err := os.Stat(f.fileDir)
	if err != nil && !os.IsExist(err) {
		os.Mkdir(f.fileDir, 0755)
	}
}

func (f *FileLogger) split() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.logFile != nil {
		f.logFile.Close()
		if f.hook != nil {
			f.hook.AfterSplit(f.logFilepath)
		}
	}

	t, _ := time.Parse(f.dateFormat, time.Now().Format(f.dateFormat))
	f.date = &t

	logFile := f.getLogFilePath()

	f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	f.lg = log.New(f.logFile, f.prefix, f.logFlag)
	return
}

// 日志写入
func (f *FileLogger) logWriter() {
	defer func() { recover() }()

	for {
		str := <-f.logChan
		f.mu.RLock()
		f.lg.Output(2, str)
		f.mu.RUnlock()
	}
}

func (f *FileLogger) Println(info string) {
	f.logChan <- info
}

func New(conf LoggerConf, names ...string) (err error) {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	multLogMut.Lock()
	defer multLogMut.Unlock()
	if _, ok := multLog[name]; ok {
		err = fmt.Errorf("日志名称:%s已经被使用", name)
		return
	}
	dateFormat := "20060102"
	if conf.DateFormat != "" {
		dateFormat = conf.DateFormat
	}
	f := &FileLogger{
		fileDir:    conf.FileDir,
		fileName:   conf.FileName,
		prefix:     conf.Prefix,
		logFlag:    conf.Flag,
		mu:         new(sync.RWMutex),
		logChan:    make(chan string, 1024),
		dateFormat: dateFormat,
		hook:       conf.Hook,
	}

	t, _ := time.Parse(dateFormat, time.Now().Format(dateFormat))
	f.date = &t

	f.isExistOrCreate()

	logFile := f.getLogFilePath()

	f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	f.lg = log.New(f.logFile, f.prefix, f.logFlag)
	go f.logWriter()
	multLog[name] = f
	return
}

func StartMonitor() {
	defer func() { recover() }()
	timer := time.NewTicker(30 * time.Second)
	for {
		<-timer.C
		multLogMut.RLock()
		for _, f := range multLog {
			if f.lg != nil && f.isMustSplit() {
				if err := f.split(); err != nil {
				}
			}
		}
		multLogMut.RUnlock()
	}
}

func Get(names ...string) *FileLogger {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	multLogMut.RLock()
	defer multLogMut.RUnlock()
	if f, ok := multLog[name]; ok {
		return f
	}
	return nil
}

// 关闭日志
func Close(names ...string) {
	var beCloseNames []string
	if len(names) == 0 {
		beCloseNames = append(beCloseNames, "")
	}
	multLogMut.Lock()
	defer multLogMut.Unlock()
	for _, name := range names {
		if f, ok := multLog[name]; ok {
			close(f.logChan)
			f.lg = nil
			f.logFile.Close()
			delete(multLog, name)
		}
	}
}

func Println(info string, names ...string) {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	multLogMut.RLock()
	if f, ok := multLog[name]; ok {
		f.logChan <- info
	} else {
		log.Printf("日志名称:%s不存在", name)
	}
	multLogMut.RUnlock()
}
