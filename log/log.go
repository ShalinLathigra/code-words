package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// Basic idea is that we're going to be constructing a JSON object as a log

const (
	MAX_LEN       int    = 512
	SUFFIX_LEN    int    = 2
	SUFFIX        string = "}\n"
	LOG_FILE_PATH string = "/home/dev/code-words/history/"
	// if when we create the first logger the files do not exist, then use os to create it
)

type LogLevel int

type LogBuilder struct {
	bytes        []byte
	length       int
	module       []byte
	moduleLength int
	file         *os.File
	writer       *bufio.Writer
}

var loggers = make(map[string]*LogBuilder)

var allLogger = CreateLogger("all")
var writeMutex = sync.Mutex{}
var startTime = time.Now()

func createWriter(key string, fileName string) (*os.File, *bufio.Writer) {
	if len(fileName)+len(key) == 0 {
		panic("Cannot init with empty log file name")
	}
	os.MkdirAll(fmt.Sprintf("%s%s", LOG_FILE_PATH, key), os.ModePerm)
	if file, err := os.OpenFile(fmt.Sprintf("%s%s/%s", LOG_FILE_PATH, key, fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		panic(err)
	} else {
		return file, bufio.NewWriter(file)
	}
}

func CreateLogger(mod string) *LogBuilder {
	if logger, ok := loggers[mod]; ok {
		return logger.Log()
	}
	file, writer := createWriter(fmt.Sprintf("%d-%d-%d", startTime.Hour(), startTime.Minute(), startTime.Second()), mod)
	loggers[mod] = &LogBuilder{
		bytes:        make([]byte, MAX_LEN),
		module:       []byte(mod),
		moduleLength: len(mod),
		file:         file,
		writer:       writer,
	}
	return loggers[mod].Log()
}

func (lb *LogBuilder) Log() *LogBuilder {
	lb.length = 0
	lb.appendBytes([]byte("{\"mod\":\""), 7)
	lb.appendBytes(lb.module, lb.moduleLength)
	lb.appendBytes([]byte("\""), 1)
	return lb
}

func (lb *LogBuilder) String(key string, val string) *LogBuilder {
	lb.appendKey(key)
	lb.appendBytes([]byte(fmt.Sprintf("\"%s\"", val)), len(val)+2)
	return lb
}

func (lb *LogBuilder) Bool(key string, val bool) *LogBuilder {
	lb.appendKey(key)
	if val {
		lb.appendBytes([]byte("true"), 4)
	} else {
		lb.appendBytes([]byte("false"), 5)
	}
	return lb
}

func (lb *LogBuilder) Int(key string, val int) *LogBuilder {
	lb.appendKey(key)
	msg := []byte(fmt.Sprintf("%d", val))
	lb.appendBytes(msg, len(msg))
	return lb
}

func (lb *LogBuilder) Any(key string, val any) *LogBuilder {
	lb.appendKey(key)
	msg := []byte(fmt.Sprintf("%v", val))
	lb.appendBytes(msg, len(msg))
	return lb
}

func (lb *LogBuilder) Msg(msg string) {
	lb.String("msg", msg)
	lb.appendBytes([]byte("}\n"), 2)
	lb.logLine(lb.bytes, lb.length)
	allLogger.logLine(lb.bytes, lb.length)
}

func (lb *LogBuilder) appendKey(key string) {
	chunk := ""
	if lb.length > 1 {
		chunk = ", "
	}
	chunk = fmt.Sprintf("%s\"%s\":", chunk, key)
	lb.appendBytes([]byte(chunk), len(chunk))

}
func (lb *LogBuilder) appendBytes(chunk []byte, len int) {
	newLen := lb.length + len + SUFFIX_LEN
	if newLen > MAX_LEN {
		panic(fmt.Errorf("log is too long, would be %d (including suffix) vs %d", newLen, MAX_LEN))
	}
	copy(lb.bytes[lb.length:lb.length+len], chunk[:])
	lb.length += len
}

func (lb *LogBuilder) logLine(chunk []byte, len int) {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	if _, err := lb.writer.Write(chunk[0:len]); err != nil {
		panic(err)
	} else {
		if err := lb.writer.Flush(); err != nil {
			panic(err)
		}
	}
}

func Clean() {
	// close files safely
	// for each logger:
	writeMutex.Lock()
	defer writeMutex.Unlock()
	for _, logger := range loggers {
		if err := logger.writer.Flush(); err != nil {
			fmt.Printf("Error clearing writer: %s", logger.module)
		}
		if err := logger.file.Close(); err != nil {
			fmt.Printf("Error closing file (%s): %s", logger.module, logger.file.Name())
		}
	}
}
