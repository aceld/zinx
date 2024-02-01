// Package zlog provides logging interfaces for zinx.
// This includes:
//
// - stdzlog module, which provides global logging methods
// - zlogger module, which defines logging protocols as object methods
//
// Current file description:
// @Title zlogger.go
// @Description Basic logging interface, including Debug, Fatal, etc.
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package zlog

/*
	All methods and APIs of the log class.
	Add By Aceld(刘丹冰) 2019-4-23
*/

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/aceld/zinx/zutils"
)

const (
	LOG_MAX_BUF = 1024 * 1024
)

// Log header information flag, using bitmap mode, users can choose which flag bits to print in the header
// (日志头部信息标记位，采用bitmap方式，用户可以选择头部需要哪些标记位被打印)
const (
	BitDate         = 1 << iota                            // Date flag bit 2019/01/23 (日期标记位)
	BitTime                                                // Time flag bit 01:23:12 (时间标记位)
	BitMicroSeconds                                        // Microsecond flag bit 01:23:12.111222 (微秒级标记位)
	BitLongFile                                            // Complete file name /home/go/src/zinx/server.go (完整文件名称)
	BitShortFile                                           // Last file name server.go (最后文件名)
	BitLevel                                               // Current log level: 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal) (当前日志级别)
	BitStdFlag      = BitDate | BitTime                    // Standard log header format (标准头部日志格式)
	BitDefault      = BitLevel | BitShortFile | BitStdFlag // Default log header format (默认日志头部格式)
)

// Log Level
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
	LogPanic
	LogFatal
)

// Log Level String
var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type ZinxLoggerCore struct {
	// to ensure thread-safe when multiple goroutines read and write files to prevent mixed-up content, achieving concurrency safety
	// (确保多协程读写文件，防止文件内容混乱，做到协程安全)
	mu sync.Mutex

	// the prefix string for each line of the log, which has the log tag
	// (每行log日志的前缀字符串,拥有日志标记)
	prefix string

	// log tag bit (日志标记位)
	flag int

	// the output buffer (输出的缓冲区)
	buf bytes.Buffer

	// log isolation level
	// (日志隔离级别)
	isolationLevel int

	// call stack depth of the function that gets the log file name and code using runtime.Call
	// (获取日志文件名和代码上述的runtime.Call 的函数调用层数)
	calldDepth int

	fw *zutils.Writer

	onLogHook func([]byte)
}

/*
NewZinxLog Create a new log

out: The file io for standard output
prefix: The prefix of the log
flag: The flag of the log header information
*/
func NewZinxLog(prefix string, flag int) *ZinxLoggerCore {

	// By default, debug is turned on, the depth is 2, and the ZinxLogger object calling the log print method can call up to two levels to reach the output function
	// (默认 debug打开， calledDepth深度为2,ZinxLogger对象调用日志打印方法最多调用两层到达output函数)
	zlog := &ZinxLoggerCore{prefix: prefix, flag: flag, isolationLevel: 0, calldDepth: 2}

	// Set the log object's resource cleanup destructor method (this is not necessary, as go's Gc will automatically collect, but for the sake of neatness)
	// (设置log对象 回收资源 析构方法(不设置也可以，go的Gc会自动回收，强迫症没办法))
	runtime.SetFinalizer(zlog, CleanZinxLog)
	return zlog
}

// CleanZinxLog Recycle log resources
func CleanZinxLog(log *ZinxLoggerCore) {
	log.closeFile()
}

func (log *ZinxLoggerCore) SetLogHook(f func([]byte)) {
	log.onLogHook = f
}

/*
formatHeader generates the header information for a log entry.

t: The current time.
file: The file name of the source code invoking the log function.
line: The line number of the source code invoking the log function.
level: The log level of the current log entry.
*/
func (log *ZinxLoggerCore) formatHeader(t time.Time, file string, line int, level int) {
	var buf *bytes.Buffer = &log.buf
	// If the current prefix string is not empty, write the prefix first.
	if log.prefix != "" {
		buf.WriteByte('<')
		buf.WriteString(log.prefix)
		buf.WriteByte('>')
	}

	// If the time-related flags are set, add the time information to the log header.
	if log.flag&(BitDate|BitTime|BitMicroSeconds) != 0 {
		// Date flag is set
		if log.flag&BitDate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/') // "2019/"
			itoa(buf, int(month), 2)
			buf.WriteByte('/') // "2019/04/"
			itoa(buf, day, 2)
			buf.WriteByte(' ') // "2019/04/11 "
		}

		// Time flag is set
		if log.flag&(BitTime|BitMicroSeconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':') // "11:"
			itoa(buf, min, 2)
			buf.WriteByte(':') // "11:15:"
			itoa(buf, sec, 2)  // "11:15:33"
			// Microsecond flag is set
			if log.flag&BitMicroSeconds != 0 {
				buf.WriteByte('.')
				itoa(buf, t.Nanosecond()/1e3, 6) // "11:15:33.123123
			}
			buf.WriteByte(' ')
		}

		// Log level flag is set
		if log.flag&BitLevel != 0 {
			buf.WriteString(levels[level])
		}

		// Short file name flag or long file name flag is set
		if log.flag&(BitShortFile|BitLongFile) != 0 {
			// Short file name flag is set
			if log.flag&BitShortFile != 0 {
				short := file
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						// Get the file name after the last '/' character, e.g. "zinx.go" from "/home/go/src/zinx.go"
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			buf.WriteString(file)
			buf.WriteByte(':')
			itoa(buf, line, -1) // line number
			buf.WriteString(": ")
		}
	}
}

// OutPut outputs log file, the original method
func (log *ZinxLoggerCore) OutPut(level int, s string) error {
	now := time.Now() // get current time
	var file string   // file name of the current caller of the log interface
	var line int      // line number of the executed code
	log.mu.Lock()
	defer log.mu.Unlock()

	if log.flag&(BitShortFile|BitLongFile) != 0 {
		log.mu.Unlock()
		var ok bool
		// get the file name and line number of the current caller
		_, file, line, ok = runtime.Caller(log.calldDepth)
		if !ok {
			file = "unknown-file"
			line = 0
		}
		log.mu.Lock()
	}

	// reset buffer
	log.buf.Reset()
	// write log header
	log.formatHeader(now, file, line, level)
	// write log content
	log.buf.WriteString(s)
	// add line break
	if len(s) > 0 && s[len(s)-1] != '\n' {
		log.buf.WriteByte('\n')
	}

	var err error
	if log.fw == nil {
		// if log file is not set, output to console
		_, _ = os.Stderr.Write(log.buf.Bytes())
	} else {
		// write the filled buffer to IO output
		_, err = log.fw.Write(log.buf.Bytes())
	}

	if log.onLogHook != nil {
		log.onLogHook(log.buf.Bytes())
	}
	return err
}

func (log *ZinxLoggerCore) verifyLogIsolation(logLevel int) bool {
	return log.isolationLevel > logLevel
}

func (log *ZinxLoggerCore) Debugf(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogDebug) {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintf(format, v...))
}

func (log *ZinxLoggerCore) Debug(v ...interface{}) {
	if log.verifyLogIsolation(LogDebug) {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintln(v...))
}

func (log *ZinxLoggerCore) Infof(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogInfo) {
		return
	}
	_ = log.OutPut(LogInfo, fmt.Sprintf(format, v...))
}

func (log *ZinxLoggerCore) Info(v ...interface{}) {
	if log.verifyLogIsolation(LogInfo) {
		return
	}
	_ = log.OutPut(LogInfo, fmt.Sprintln(v...))
}

func (log *ZinxLoggerCore) Warnf(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogWarn) {
		return
	}
	_ = log.OutPut(LogWarn, fmt.Sprintf(format, v...))
}

func (log *ZinxLoggerCore) Warn(v ...interface{}) {
	if log.verifyLogIsolation(LogWarn) {
		return
	}
	_ = log.OutPut(LogWarn, fmt.Sprintln(v...))
}

func (log *ZinxLoggerCore) Errorf(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogError) {
		return
	}
	_ = log.OutPut(LogError, fmt.Sprintf(format, v...))
}

func (log *ZinxLoggerCore) Error(v ...interface{}) {
	if log.verifyLogIsolation(LogError) {
		return
	}
	_ = log.OutPut(LogError, fmt.Sprintln(v...))
}

func (log *ZinxLoggerCore) Fatalf(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogFatal) {
		return
	}
	_ = log.OutPut(LogFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (log *ZinxLoggerCore) Fatal(v ...interface{}) {
	if log.verifyLogIsolation(LogFatal) {
		return
	}
	_ = log.OutPut(LogFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

func (log *ZinxLoggerCore) Panicf(format string, v ...interface{}) {
	if log.verifyLogIsolation(LogPanic) {
		return
	}
	s := fmt.Sprintf(format, v...)
	_ = log.OutPut(LogPanic, s)
	panic(s)
}

func (log *ZinxLoggerCore) Panic(v ...interface{}) {
	if log.verifyLogIsolation(LogPanic) {
		return
	}
	s := fmt.Sprintln(v...)
	_ = log.OutPut(LogPanic, s)
	panic(s)
}

func (log *ZinxLoggerCore) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, LOG_MAX_BUF)
	n := runtime.Stack(buf, true) //得到当前堆栈信息
	s += string(buf[:n])
	s += "\n"
	_ = log.OutPut(LogError, s)
}

// Flags gets the current log bitmap flags
// (获取当前日志bitmap标记)
func (log *ZinxLoggerCore) Flags() int {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.flag
}

// ResetFlags resets the log Flags bitmap flags
// (重新设置日志Flags bitMap 标记位)
func (log *ZinxLoggerCore) ResetFlags(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.flag = flag
}

// AddFlag adds a flag to the bitmap flags
// (添加flag标记)
func (log *ZinxLoggerCore) AddFlag(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.flag |= flag
}

// SetPrefix sets a custom prefix for the log
// (设置日志的 用户自定义前缀字符串)
func (log *ZinxLoggerCore) SetPrefix(prefix string) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.prefix = prefix
}

// SetLogFile sets the log file output
// (设置日志文件输出)
func (log *ZinxLoggerCore) SetLogFile(fileDir string, fileName string) {
	if log.fw != nil {
		log.fw.Close()
	}
	log.fw = zutils.New(filepath.Join(fileDir, fileName))
}

// SetMaxAge 最大保留天数
func (log *ZinxLoggerCore) SetMaxAge(ma int) {
	if log.fw == nil {
		return
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	log.fw.SetMaxAge(ma)
}

// SetMaxSize 单个日志最大容量 单位：字节
func (log *ZinxLoggerCore) SetMaxSize(ms int64) {
	if log.fw == nil {
		return
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	log.fw.SetMaxSize(ms)
}

// SetCons 同时输出控制台
func (log *ZinxLoggerCore) SetCons(b bool) {
	if log.fw == nil {
		return
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	log.fw.SetCons(b)
}

// Close the file associated with the log
// (关闭日志绑定的文件)
func (log *ZinxLoggerCore) closeFile() {
	if log.fw != nil {
		log.fw.Close()
	}
}

func (log *ZinxLoggerCore) SetLogLevel(logLevel int) {
	log.isolationLevel = logLevel
}

// Convert an integer to a fixed-length string, where the width of the string should be greater than 0
// Ensure that the buffer has sufficient capacity
// (将一个整形转换成一个固定长度的字符串，字符串宽度应该是大于0的
// 要确保buffer是有容量空间的)
func itoa(buf *bytes.Buffer, i int, wID int) {
	var u uint = uint(i)
	if u == 0 && wID <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wID > 0; u /= 10 {
		bp--
		wID--
		b[bp] = byte(u%10) + '0'
	}

	// avoID slicing b to avoID an allocation.
	for bp < len(b) {
		buf.WriteByte(b[bp])
		bp++
	}
}
