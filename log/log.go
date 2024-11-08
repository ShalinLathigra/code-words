package log

import (
	"bufio"
	"fmt"
	"os"
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

var startTime = time.Now()

func createWriter(key string, fileName string) (*os.File, *bufio.Writer) {
	if len(fileName)+len(key) == 0 {
		panic("Cannot init with empty log file name")
	}
	os.MkdirAll(fmt.Sprintf("%s%s", LOG_FILE_PATH, key), os.ModePerm)
	if file, err := os.Create(fmt.Sprintf("%s%s/%s", LOG_FILE_PATH, key, fileName)); err != nil {
		panic(err)
	} else {
		return file, bufio.NewWriter(file)
	}
}

func CreateLogger(mod string) *LogBuilder {
	file, writer := createWriter(fmt.Sprintf("%d-%d-%d", startTime.Hour(), startTime.Minute(), startTime.Second()), mod)
	loggers[mod] = &LogBuilder{
		make([]byte, MAX_LEN),
		0,
		[]byte(mod),
		len(mod),
		file,
		writer,
	}
	return loggers[mod].Log()
}

func (lb *LogBuilder) Log() *LogBuilder {
	lb.length = 0
	lb.appendBytes([]byte("{\"mod\":\""), 9)
	lb.appendBytes(lb.module, lb.moduleLength)
	lb.appendBytes([]byte("\""), 1)
	return lb
}

func (lb *LogBuilder) String(key string, val string) *LogBuilder {
	lb.appendKey(key)
	lb.appendBytes([]byte(fmt.Sprintf("\"%s\"", val)), len(val)+2)
	return lb
}

func (lb *LogBuilder) Msg(msg string) {
	lb.String("msg", msg)
	lb.appendBytes([]byte("}\n"), 2)
	lb.length += 2
	lb.logLine(lb.bytes, lb.length)
	lb.length = 0
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
	if nn, err := lb.writer.Write(chunk[0:len]); err != nil {
		panic(err)
	} else {
		fmt.Printf("Writing (%d) bytes: (%d) %s", nn, len, chunk)
		if err := lb.writer.Flush(); err != nil {
			panic(err)
		}
	}
}

func Clean() {
	// close files safely
	// for each logger:
	for _, logger := range loggers {
		if err := logger.writer.Flush(); err != nil {
			fmt.Printf("Error clearing writer: %s", logger.module)
		}
		if err := logger.file.Close(); err != nil {
			fmt.Printf("Error closing file (%s): %s", logger.module, logger.file.Name())
		}
	}
}
